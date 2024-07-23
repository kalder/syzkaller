// Copyright 2024 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package coveragedb

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/google/syzkaller/pkg/subsystem"
	_ "github.com/google/syzkaller/pkg/subsystem/lists"
	"github.com/google/uuid"
)

type FilesRecord struct {
	Session      string
	FilePath     string
	Instrumented int64
	Covered      int64
}

type FileSubsystems struct {
	Namespace  string
	FilePath   string
	Subsystems []string
}

type HistoryRecord struct {
	Session   string
	Time      time.Time
	Namespace string
	Repo      string
	Commit    string
	Duration  int64
	DateTo    civil.Date
	TotalRows int64
}

func NewClient(ctx context.Context, projectID string) (*spanner.Client, error) {
	database := "projects/" + projectID + "/instances/syzbot/databases/coverage"
	return spanner.NewClient(ctx, database)
}

// TODO: move to dashAPI once tested? I'm not sure we'll benefit.

type Coverage struct {
	Instrumented int64
	Covered      int64
}

func SaveMergeResult(ctx context.Context, projectID string, covMap map[string]*Coverage,
	template *HistoryRecord, totalRows int64, sss []*subsystem.Subsystem) error {
	client, err := NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("spanner.NewClient() failed: %s", err.Error())
	}
	defer client.Close()

	ssMatcher := subsystem.MakePathMatcher(sss)
	ssCache := make(map[string][]string)

	session := uuid.New().String()
	mutations := []*spanner.Mutation{}
	for filePath, record := range covMap {
		mutations = append(mutations, fileRecordMutation(session, filePath, record))
		subsystems := fileSubsystems(filePath, ssMatcher, ssCache)
		mutations = append(mutations, fileSubsystemsMutation(template.Namespace, filePath, subsystems))
		// 80k mutations is a DB limit. 4 fields * 2k records is apx 8k mutations
		// let keep this value 10x lower to have a room for indexes
		// indexes update are also counted
		if len(mutations) > 2000 {
			if _, err = client.Apply(ctx, mutations); err != nil {
				return fmt.Errorf("failed to spanner.Apply(inserts): %s", err.Error())
			}
			mutations = nil
		}
	}
	mutations = append(mutations, historyMutation(session, template, totalRows))
	if _, err = client.Apply(ctx, mutations); err != nil {
		return fmt.Errorf("failed to spanner.Apply(inserts): %s", err.Error())
	}
	return nil
}

func historyMutation(session string, template *HistoryRecord, totalRows int64) *spanner.Mutation {
	historyInsert, err := spanner.InsertOrUpdateStruct("merge_history", &HistoryRecord{
		Session:   session,
		Time:      time.Now(),
		Namespace: template.Namespace,
		Repo:      template.Repo,
		Commit:    template.Commit,
		Duration:  template.Duration,
		DateTo:    template.DateTo,
		TotalRows: totalRows,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to spanner.InsertStruct(): %s", err.Error()))
	}
	return historyInsert
}

func fileRecordMutation(session, filePath string, record *Coverage) *spanner.Mutation {
	insert, err := spanner.InsertOrUpdateStruct("files", &FilesRecord{
		Session:      session,
		FilePath:     filePath,
		Instrumented: record.Instrumented,
		Covered:      record.Covered,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to fileRecordMutation(): %s", err.Error()))
	}
	return insert
}

func fileSubsystemsMutation(ns, filePath string, subsystems []string) *spanner.Mutation {
	insert, err := spanner.InsertOrUpdateStruct("file_subsystems", &FileSubsystems{
		Namespace:  ns,
		FilePath:   filePath,
		Subsystems: subsystems,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to fileSubsystemsMutation(): %s", err.Error()))
	}
	return insert
}

func fileSubsystems(filePath string, ssMatcher *subsystem.PathMatcher, ssCache map[string][]string) []string {
	sss, cached := ssCache[filePath]
	if !cached {
		for _, match := range ssMatcher.Match(filePath) {
			sss = append(sss, match.Name)
		}
		ssCache[filePath] = sss
	}
	return sss
}
