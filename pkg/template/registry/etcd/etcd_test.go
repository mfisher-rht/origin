package etcd

import (
	"testing"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/registry/registrytest"
	"k8s.io/kubernetes/pkg/storage"
	etcdstorage "k8s.io/kubernetes/pkg/storage/etcd"
	"k8s.io/kubernetes/pkg/tools"
	"k8s.io/kubernetes/pkg/tools/etcdtest"

	"github.com/openshift/origin/pkg/api/latest"
	"github.com/openshift/origin/pkg/template/api"
)

func newHelper(t *testing.T) (*tools.FakeEtcdClient, storage.Interface) {
	fakeEtcdClient := tools.NewFakeEtcdClient(t)
	fakeEtcdClient.TestIndex = true
	helper := etcdstorage.NewEtcdStorage(fakeEtcdClient, latest.Codec, etcdtest.PathPrefix())
	return fakeEtcdClient, helper
}

func newStorage(t *testing.T) (*REST, *tools.FakeEtcdClient) {
	etcdStorage, fakeClient := registrytest.NewEtcdStorage(t, "")
	return NewREST(etcdStorage), fakeClient
}

func validNew() *api.Template {
	return &api.Template{
		ObjectMeta: kapi.ObjectMeta{
			Name:      "foo",
			Namespace: kapi.NamespaceDefault,
		},
	}
}

func validChanged() *api.Template {
	template := validNew()
	template.ResourceVersion = "1"
	template.Labels = map[string]string{
		"foo": "bar",
	}
	return template
}

func TestStorage(t *testing.T) {
	storage, _ := newStorage(t)
	var _ rest.Creater = storage
	var _ rest.Lister = storage
	var _ rest.GracefulDeleter = storage
	var _ rest.Updater = storage
	var _ rest.Getter = storage
}

func TestCreate(t *testing.T) {
	storage, fakeClient := newStorage(t)
	test := registrytest.New(t, fakeClient, storage.Etcd)
	template := validNew()
	template.ObjectMeta = kapi.ObjectMeta{}
	test.TestCreate(
		// valid
		template,
		// invalid
		&api.Template{},
	)
}
