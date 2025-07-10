package etsi119612_test

import (
	"testing"
	"time"
	"context"
	"io"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/SUNET/g119612/pkg/etsi119612/cache"

)
var (
	customCache *cache.CustomCache
	getUrl    ="https://se-tl.se"
	data  []byte
)

func TestSetCache(t *testing.T){
	cntx :=context.Background()
	xmlFile, err := os.Open("testdata/SE-TL.xml")
	assert.NoError(t, err)
	defer xmlFile.Close()
	data, err =io.ReadAll(xmlFile)
	assert.NoError(t, err)
	customCache = cache.NewCustomCache()
	customCache.Set(cntx, getUrl, data, 60*time.Second)
	entry, ok := customCache.Data[getUrl]
	assert.True(t, ok)
	assert.Equal(t, data, entry.Value)
	expectedTime :=time.Now().Add(60*time.Second)
	assert.WithinDuration(t, expectedTime,entry.ExpiresAt, time.Second)
}

func TestGetCacheSuccess(t *testing.T){
	cntx :=context.Background()
	assert.NotNil(t, customCache)
	entry, exists :=customCache.Get(cntx, getUrl)
	assert.True(t,exists)
	assert.Equal(t, data, entry.Value)
	expectedTime :=time.Now().Add(60*time.Second)
	assert.WithinDuration(t, expectedTime,entry.ExpiresAt, time.Second)
}

func TestDeleteSuccess(t *testing.T){
	cntx :=context.Background()
	assert.NotNil(t, customCache)
	entry, exists :=customCache.Get(cntx, getUrl)
	assert.True(t,exists)
	assert.Equal(t, data, entry.Value)
	expectedTime :=time.Now().Add(60*time.Second)
	assert.WithinDuration(t, expectedTime,entry.ExpiresAt, time.Second)
	err :=customCache.Delete(cntx, getUrl)
	assert.NoError(t, err)
	entry, exists=customCache.Get(cntx, getUrl)
	assert.False(t,exists)
}

func TestDeleteError(t *testing.T){
	cntx :=context.Background()
	err :=customCache.Delete(cntx, getUrl)
	assert.Error(t, err)
}
