package bolt

import (
	lang "github.com/Aptomi/aptomi/pkg/slinga/language"
	"github.com/Aptomi/aptomi/pkg/slinga/object"
	"github.com/Aptomi/aptomi/pkg/slinga/object/codec/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestBoltStore(t *testing.T) {
	catalog := object.NewObjectCatalog(lang.ServiceObject, lang.ContractObject, lang.ClusterObject, lang.RuleObject, lang.DependencyObject)
	db := NewBoltStore(catalog, yaml.NewCodec(catalog))

	f, err := ioutil.TempFile("", t.Name())
	assert.Nil(t, err, "Temp file should be successfully created")
	defer os.Remove(f.Name())

	err = db.Open(f.Name())
	if err != nil {
		panic(err)
	}

	policy := lang.LoadUnitTestsPolicy("../../../testdata/unittests")

	services := make([]object.Base, 0, len(policy.Services))

	for _, service := range policy.Services {
		updated, err := db.Save(service)
		if err != nil {
			panic(err)
		}
		services = append(services, service)
		assert.True(t, updated, "Object saved for the first time")
	}

	assert.Equal(t, 2, len(services), "Len!")

	for _, service := range services {
		obj, err := db.GetByName(service.GetNamespace(), service.GetKind(), service.GetName(), service.GetGeneration())
		if err != nil {
			panic(err)
		}

		assert.Exactly(t, service, obj, "fail!")
	}
}
