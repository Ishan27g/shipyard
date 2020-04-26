package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/shipyard-run/shipyard/pkg/clients/mocks"
	"github.com/shipyard-run/shipyard/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPurgeCommand(t *testing.T) (*cobra.Command, *mocks.MockDocker, *mocks.ImageLog, func()) {
	home := os.Getenv("HOME")

	// create a fake home folder
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("HOME", dir)

	// create the fake blueprints
	err = os.MkdirAll(utils.GetBlueprintLocalFolder(""), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// create the fake helm
	err = os.MkdirAll(utils.GetHelmLocalFolder(""), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	mockDocker := &mocks.MockDocker{}
	mockDocker.On("ImageRemove", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockDocker.On("VolumeRemove", mock.Anything, mock.Anything, true).Return(nil)

	mockImageLog := &mocks.ImageLog{}
	mockImageLog.On("Read", mock.Anything).Return([]string{"one", "two"}, nil)
	mockImageLog.On("Clear").Return(nil)

	pc := newPurgeCmd(mockDocker, mockImageLog, hclog.NewNullLogger())

	return pc, mockDocker, mockImageLog, func() {
		os.RemoveAll(dir)
		os.Setenv("HOME", home)
	}
}

func TestPurgeCallsImageRemoveForCachedImages(t *testing.T) {
	pc, md, mi, cleanup := setupPurgeCommand(t)
	defer cleanup()

	err := pc.Execute()

	assert.NoError(t, err)
	md.AssertNumberOfCalls(t, "ImageRemove", 2)
	mi.AssertCalled(t, "Clear")
}

func TestPurgeReturnsErrorOnImageRemoveError(t *testing.T) {
	pc, md, mi, cleanup := setupPurgeCommand(t)
	defer cleanup()
	removeOn(&md.Mock, "ImageRemove")
	md.On("ImageRemove", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Boom"))

	err := pc.Execute()

	assert.Error(t, err)
	md.AssertNumberOfCalls(t, "ImageRemove", 1)
	mi.AssertNotCalled(t, "Clear")
}

func TestPurgeRemovesBlueprints(t *testing.T) {
	pc, _, _, cleanup := setupPurgeCommand(t)
	defer cleanup()

	err := pc.Execute()

	assert.NoError(t, err)
	assert.NoDirExists(t, utils.GetBlueprintLocalFolder(""))
}

func TestPurgeRemovesHelmCharts(t *testing.T) {
	pc, _, _, cleanup := setupPurgeCommand(t)
	defer cleanup()

	err := pc.Execute()

	assert.NoError(t, err)
	assert.NoDirExists(t, utils.GetHelmLocalFolder(""))
}
