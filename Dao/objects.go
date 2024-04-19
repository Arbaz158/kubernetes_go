package Dao

import (
	"KubernetesGo/Models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/cdedev/response"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	av1 "k8s.io/api/rbac/v1"
	v11 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var clientset *kubernetes.Clientset
var clientConfig *rest.Config
var connectionString = "root:root@tcp(127.0.0.1:3306)/kubernetes_cluster_credentials?parseTime=true"

type UserImpl struct {
}

//Checking Cluster Credentials from the passed Information

func GettingClusterCred(data Models.Credentials, w http.ResponseWriter, r *http.Request) *kubernetes.Clientset {
	var err error
	config := api.NewConfig()
	config.Kind = "Config"
	config.APIVersion = "v1"
	config.Clusters[data.Cluster] = &api.Cluster{
		Server:                   data.Server,
		CertificateAuthorityData: data.CertificateAuthorityData,
	}
	config.AuthInfos[data.Cluster] = &api.AuthInfo{
		ClientCertificateData: data.ClientCertificateData,
		ClientKeyData:         data.ClientKeyData,
		Token:                 data.Token,
	}
	config.Contexts[data.Cluster] = &api.Context{
		Cluster:  data.Cluster,
		AuthInfo: data.Cluster,
	}
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*config, data.Cluster, &clientcmd.ConfigOverrides{}, nil)
	clientConfig, err = clientBuilder.ClientConfig()
	if err != nil {
		logrus.Error("ResponseCode:", 2029, "Results:", err)
		response.RespondJSON(w, 2029, nil)
	}
	clientset, err = kubernetes.NewForConfig(clientConfig)
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, nil)
	}
	return clientset
}

func CheckRemoteConnection(data Models.Credentials, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var client *kubernetes.Clientset
	var err error
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, "Check your json values")
	}
	client = nil
	config := api.NewConfig()
	config.Kind = "Config"
	config.APIVersion = "v1"
	config.Clusters[data.Cluster] = &api.Cluster{
		Server:                   data.Server,
		CertificateAuthorityData: data.CertificateAuthorityData,
	}
	config.AuthInfos[data.Cluster] = &api.AuthInfo{
		ClientCertificateData: data.ClientCertificateData,
		ClientKeyData:         data.ClientKeyData,
		Token:                 data.Token,
	}
	config.Contexts[data.Cluster] = &api.Context{
		Cluster:  data.Cluster,
		AuthInfo: data.Cluster,
	}
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*config, data.Cluster, &clientcmd.ConfigOverrides{}, nil)
	clientConfig, err = clientBuilder.ClientConfig()
	if err != nil {
		logrus.Error("ResponseCode:", 2029, "Results:", err)
		response.RespondJSON(w, 2029, nil)
	}
	client, _ = kubernetes.NewForConfig(clientConfig)
	_, nerr := client.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
	if nerr != nil {
		logrus.Error("ResponseCode:", 2090, "Results:", nerr)
		response.RespondJSON(w, 2090, nil)
	} else {
		DatabaseConnect(data, w, r)
	}
}

func DatabaseConnect(data Models.Credentials, w http.ResponseWriter, r *http.Request) {
	var err error
	Database, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		logrus.Error("ResponseCode:", 2091, "Results:", err)
		response.RespondJSON(w, 2091, nil)
	}
	err = Database.AutoMigrate(&Models.Credentials{})
	if err != nil {
		logrus.Error("ResponseCode:", 2024, "Results:", err)
		response.RespondJSON(w, 2024, nil)
	}
	err = Database.Create(&data).Error
	if err != nil {
		logrus.Error("ResponseCode:", 2092, "Results:", err)
		response.RespondJSON(w, 2092, data.ClusterId)
	} else {
		response.RespondJSON(w, 1000, nil)
	}

}

func GettingCredentialsBasedOnClusterId(cluster_id string, w http.ResponseWriter, r *http.Request) {
	var res Models.Credentials
	Database, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		logrus.Error("ResponseCode:", 2091, "Results:", err)
		response.RespondJSON(w, 2091, nil)
	}
	err = Database.Where("cluster_id = ?", cluster_id).Find(&res).Error
	if err != nil {
		logrus.Error("ResponseCode:", 2020, "Results:", err)
		response.RespondJSON(w, 2020, cluster_id)
	}
	GettingClusterCred(res, w, r)
}

//Getting details about nodes

func (u UserImpl) ListClusterNodes(w http.ResponseWriter, r *http.Request, cluster_id string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetClusterNodes(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

//creating Namespace

func (u UserImpl) CreateNS(n Models.Namespace, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, nil)
	}
	GettingCredentialsBasedOnClusterId(n.ClusterId, w, r)
	var labels map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(n.Labels); i++ {
		labels[n.Labels[i].FirstLabel] = n.Labels[i].SecondLabel
	}

	create := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   n.Name,
			Labels: labels,
		},
	}
	result, err := clientset.CoreV1().Namespaces().Create(context.TODO(), create, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, result)
	}
}

func (u UserImpl) GetNs(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

//Getting List of All namespaces

func (u UserImpl) ListNS(w http.ResponseWriter, r *http.Request, cluster_id string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}

}

// //Deleting Namespace

func (u UserImpl) DeleteNS(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// //Update Namespace

// func (u UserImpl) UpdateNS(update Models.Namespace) (*v1.Namespace, error) {
// 	GettingCredentialsBasedOnClusterId(update.ClusterId)
// 	var labels map[string]string
// 	labels = make(map[string]string)
// 	for i := 0; i < len(update.Labels); i++ {
// 		labels[update.Labels[i].FirstLabel] = update.Labels[i].SecondLabel
// 	}

// 	nsUpdate := &v1.Namespace{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:   update.Name,
// 			Labels: labels,
// 		},
// 	}
// 	result, err := clientset.CoreV1().Namespaces().Update(context.TODO(), nsUpdate, metav1.UpdateOptions{})
// 	if err != nil {
// 		fmt.Println("This error occurs inside Dao update function", err)
// 	}
// 	return result, err
// }

//Creating Service Account

func (u UserImpl) CreateSA(w http.ResponseWriter, r *http.Request, service Models.ServiceAccount) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(service.ClusterId, w, r)
	var labels map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(service.Labels); i++ {
		labels[service.Labels[i].FirstLabel] = service.Labels[i].SecondLabel
	}

	saClient := clientset.CoreV1().ServiceAccounts(v1.NamespaceDefault)
	saCreate := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels:    labels,
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			{
				Name: service.SecretName,
			},
		},
	}

	//Creating Service Account
	result, err := saClient.Create(context.TODO(), saCreate, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, result)
	}

}

// //Getting Service Account

func (u UserImpl) GetSA(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	saGet, err := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, saGet)
	}
}

// //Listing Service Account

func (u UserImpl) ListSA(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	saList, err := clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, saList)
	}

}

//Deleting Service account

func (u UserImpl) DeleteSA(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}

}

// //Updating Service Account

// func (u UserImpl) UpdateSA(update Models.ServiceAccount, namespace string) (*v1.ServiceAccount, error) {
// 	GettingCredentialsBasedOnClusterId(update.ClusterId)
// 	var labels map[string]string
// 	labels = make(map[string]string)
// 	for i := 0; i < len(update.Labels); i++ {
// 		labels[update.Labels[i].FirstLabel] = update.Labels[i].SecondLabel
// 	}

// 	saUpdate := &v1.ServiceAccount{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      update.Name,
// 			Namespace: update.Namespace,
// 			Labels:    labels,
// 		},
// 		ImagePullSecrets: []v1.LocalObjectReference{
// 			{
// 				Name: update.SecretName,
// 			},
// 		},
// 	}
// 	sa, err := clientset.CoreV1().ServiceAccounts(namespace).Update(context.TODO(), saUpdate, metav1.UpdateOptions{})
// 	if err != nil {
// 		fmt.Println("This error occurs inside Update Dao function", err)
// 	}
// 	return sa, err
// }

// Configmaps

func (u UserImpl) CreateConfigmap(w http.ResponseWriter, r *http.Request, namespace string, cfg Models.Configmap) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(cfg.ClusterId, w, r)
	var labels map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(cfg.Labels); i++ {
		labels[cfg.Labels[i].FirstLabel] = cfg.Labels[i].SecondLabel
	}

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
			Labels:    labels,
		},
		Data: cfg.Data,
	}

	create, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, create)
	}

}

func (u UserImpl) GetConfigmap(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) ListConfigmap(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) DeleteConfigmap(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// StorageClass

func (u UserImpl) CreateStorageClass(w http.ResponseWriter, r *http.Request, sc Models.StorageClass) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&sc)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(sc.ClusterId, w, r)
	var labels map[string]string
	var parameters map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(sc.Labels); i++ {
		labels[sc.Labels[i].FirstLabel] = sc.Labels[i].SecondLabel
	}

	parameters = make(map[string]string)
	for j := 0; j < len(sc.Parameters); j++ {
		parameters[sc.Parameters[j].FirstParameter] = sc.Parameters[j].SecondParameter
	}

	strg := &v11.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sc.Name,
			Namespace: sc.Namespace,
			Labels:    labels,
		},
		Provisioner:          sc.Provisioner,
		Parameters:           parameters,
		VolumeBindingMode:    (*v11.VolumeBindingMode)(&sc.VolumeBinidingMode),
		ReclaimPolicy:        (*v1.PersistentVolumeReclaimPolicy)(&sc.ReclaimPolicy),
		AllowVolumeExpansion: &sc.VolumeExpansion,
	}

	create, err := clientset.StorageV1().StorageClasses().Create(context.TODO(), strg, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, create)
	}
}

func (u UserImpl) GetStorageClassDetails(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) ListSC(w http.ResponseWriter, r *http.Request, cluster_id string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) DeleteSC(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.StorageV1().StorageClasses().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

//Creating Persistent Volume

func (u UserImpl) CreatePV(w http.ResponseWriter, r *http.Request, pv Models.PersistentVolume) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&pv)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(pv.ClusterId, w, r)
	var labels map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(pv.Labels); i++ {
		labels[pv.Labels[i].FirstLabel] = pv.Labels[i].SecondLabel
	}

	if pv.PVSource == "HostPath" || pv.PVSource == "hostPath" || pv.PVSource == "hostpath" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},
				NodeAffinity: &v1.VolumeNodeAffinity{
					Required: &v1.NodeSelector{
						NodeSelectorTerms: []v1.NodeSelectorTerm{
							{
								MatchExpressions: []v1.NodeSelectorRequirement{
									{
										Key:      pv.NodeKey,
										Operator: v1.NodeSelectorOperator(pv.NodeOperator),
										Values:   pv.NodeValues,
									},
								},
							},
						},
					},
				},
				PersistentVolumeSource: v1.PersistentVolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: pv.HostPath.Path,
						Type: (*v1.HostPathType)(&pv.HostPath.Type),
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}

	} else if pv.PVSource == "NFS" || pv.PVSource == "nfs" || pv.PVSource == "Nfs" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},

				PersistentVolumeSource: v1.PersistentVolumeSource{
					NFS: &v1.NFSVolumeSource{
						Server: pv.NFS.Server,
						Path:   pv.NFS.Path,
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}

	} else if pv.PVSource == "AWSEBS" || pv.PVSource == "EBS" || pv.PVSource == "AwsEbs" || pv.PVSource == "Ebs" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},
				PersistentVolumeSource: v1.PersistentVolumeSource{
					AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{
						VolumeID: pv.AwsEBS.VolumeId,
						ReadOnly: pv.AwsEBS.ReadOnly,
						FSType:   pv.AwsEBS.FsType,
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}

	} else if pv.PVSource == "AzureDisk" || pv.PVSource == "AzureDISK" || pv.PVSource == "AZDisk" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},
				PersistentVolumeSource: v1.PersistentVolumeSource{
					CSI: &v1.CSIPersistentVolumeSource{
						Driver:       pv.AzureDisk.Driver,
						VolumeHandle: pv.AzureDisk.VolumeHandle,
						ReadOnly:     pv.AzureDisk.ReadOnly,
						FSType:       pv.AzureDisk.FsType,
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}
	} else if pv.PVSource == "AzureFile" || pv.PVSource == "AZFile" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},
				MountOptions: pv.MountOptions,
				PersistentVolumeSource: v1.PersistentVolumeSource{
					CSI: &v1.CSIPersistentVolumeSource{
						Driver:       pv.AzureFile.Driver,
						VolumeHandle: pv.AzureFile.VolumeHandle,
						ReadOnly:     pv.AzureFile.ReadOnly,
						FSType:       pv.AzureFile.FsType,
						NodeStageSecretRef: &v1.SecretReference{
							Name:      pv.AzureFile.SecretName,
							Namespace: pv.AzureFile.SecretName,
						},
						VolumeAttributes: map[string]string{
							pv.AzureFile.FirstAttribute: pv.AzureFile.SecondAttribute,
						},
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}

	} else if pv.PVSource == "GKEPd" || pv.PVSource == "Gkepd" {
		pvCreate := &v1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.Name,
				Namespace: pv.Namespace,
				Labels:    labels,
			},
			Spec: v1.PersistentVolumeSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.PersistentVolumeAccessMode(pv.AccessModes),
				},
				StorageClassName: pv.StorageClass,
				Capacity: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pv.Capacity),
				},
				MountOptions: pv.MountOptions,
				PersistentVolumeSource: v1.PersistentVolumeSource{
					CSI: &v1.CSIPersistentVolumeSource{
						Driver:       pv.GkePD.Driver,
						VolumeHandle: pv.GkePD.VolumeHandle,
						ReadOnly:     pv.GkePD.ReadOnly,
						FSType:       pv.GkePD.FsType,
						VolumeAttributes: map[string]string{
							pv.GkePD.FirstAttribute: pv.GkePD.SecondAttribute,
						},
					},
				},
			},
		}
		result, err := clientset.CoreV1().PersistentVolumes().Create(context.TODO(), pvCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, result)
		}
	}
}

func (u UserImpl) ListPV(w http.ResponseWriter, r *http.Request, cluster_id string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetPV(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeletePV(w http.ResponseWriter, r *http.Request, cluster_id string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

//Persistent Volume Claim

func (u UserImpl) CreatePVC(w http.ResponseWriter, r *http.Request, namespace string, pvc Models.PersistentVolumeClaim) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&pvc)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(pvc.ClusterId, w, r)
	var labels map[string]string
	labels = make(map[string]string)
	for i := 0; i < len(pvc.Labels); i++ {
		labels[pvc.Labels[i].FirstLabel] = pvc.Labels[i].SecondLabel
	}

	persist := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvc.Name,
			Namespace: pvc.Namespace,
			Labels:    labels,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.PersistentVolumeAccessMode(pvc.AccessModes),
			},
			StorageClassName: pvc.StorageClass,
			Resources: v1.ResourceRequirements{
				Requests: map[v1.ResourceName]resource.Quantity{
					v1.ResourceStorage: resource.MustParse(pvc.Capacity),
				},
			},
			VolumeName: pvc.VolumeName,
		},
	}

	create, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), persist, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, create)
	}
}

func (u UserImpl) ListPVC(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetPVC(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "applcation/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeletePVC(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "applcation/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// Services

func (u UserImpl) CreateServices(w http.ResponseWriter, r *http.Request, namespace string, svc Models.Services) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&svc)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(svc.ClusterId, w, r)
	var labels map[string]string
	var ports []v1.ServicePort

	labels = make(map[string]string)
	for i := 0; i < len(svc.Labels); i++ {
		labels[svc.Labels[i].FirstLabel] = svc.Labels[i].SecondLabel
	}

	for i := 0; i < len(svc.NodePort); i++ {
		portsArr := []v1.ServicePort{
			{
				Name:       svc.NodePort[i].PortName,
				Port:       svc.NodePort[i].Port,
				Protocol:   v1.Protocol(svc.NodePort[i].Protocol),
				TargetPort: svc.NodePort[i].TargetPort,
				NodePort:   svc.NodePort[i].NodePort,
			},
		}
		ports = append(ports, portsArr...)

	}

	if svc.Type == "NodePort" || svc.Type == "nodeport" || svc.Type == "Nodeport" {
		svcCreate := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Labels:    labels,
			},
			Spec: v1.ServiceSpec{
				Type:     v1.ServiceType(svc.Type),
				Selector: labels,
				Ports:    ports,
			},
		}
		create, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), svcCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}

	} else if svc.Type == "ClusterIp" || svc.Type == "ClusterIP" || svc.Type == "clusterip" {
		svcCreate := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Labels:    labels,
			},
			Spec: v1.ServiceSpec{
				Type:     v1.ServiceType(svc.Type),
				Selector: labels,
				Ports: []v1.ServicePort{
					{
						Port:       svc.ClusterIp.Port,
						Name:       svc.ClusterIp.PortName,
						Protocol:   v1.Protocol(svc.ClusterIp.Protocol),
						TargetPort: svc.ClusterIp.TargetPort,
					},
				},
			},
		}
		create, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), svcCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	}
}

func (u UserImpl) ListServices(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetServices(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeleteServices(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

//Secrets

func (u UserImpl) CreateSecret(w http.ResponseWriter, r *http.Request, namespace string, secret Models.Secret) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&secret)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(secret.ClusterId, w, r)
	var labels map[string]string
	var opaque map[string][]byte
	labels = make(map[string]string)
	for i := 0; i < len(secret.Labels); i++ {
		labels[secret.Labels[i].FirstLabel] = secret.Labels[i].SecondLabel
	}

	opaque = make(map[string][]byte)
	for j := 0; j < len(secret.OpaqueSecret); j++ {
		opaque[secret.OpaqueSecret[j].FirstSecret] = secret.OpaqueSecret[j].SecondSecret
	}

	if secret.Type == "kubernetes.io/dockerconfigjson" {
		createSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				Labels:    labels,
			},
			Type: v1.SecretType(secret.Type),
			Data: map[string][]byte{
				v1.DockerConfigJsonKey: []byte(secret.DockerConfig),
			},
		}
		create, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), createSecret, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}

	} else if secret.Type == "kubernetes.io/service-account-token" {
		createSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				Labels:    labels,
			},
			Type: v1.SecretType(secret.Type),
			Data: map[string][]byte{},
		}
		create, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), createSecret, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if secret.Type == "opaque" {
		createSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				Labels:    labels,
			},
			Type: v1.SecretType(secret.Type),
			Data: opaque,
		}
		create, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), createSecret, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	}
}

func (u UserImpl) ListSecrets(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetSecret(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeleteSecret(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// Pods

func (u UserImpl) CreatePod(w http.ResponseWriter, r *http.Request, namespace string, pods Models.Pod) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&pods)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(pods.ClusterId, w, r)
	var labels map[string]string
	var ports []v1.ContainerPort
	var mountPath []v1.VolumeMount
	var env []v1.EnvVar
	var secretEnv []v1.EnvVar
	var configEnv []v1.EnvVar
	var cnt []v1.Container
	var initContainers []v1.Container
	var security v1.PodSecurityContext
	var cntsecurity v1.SecurityContext
	var readiness *v1.Probe
	var liveness *v1.Probe
	var resources v1.ResourceRequirements

	labels = make(map[string]string)
	for i := 0; i < len(pods.Labels); i++ {
		labels[pods.Labels[i].FirstLabel] = pods.Labels[i].SecondLabel
	}

	for j := 0; j < len(pods.Container.ContainerPort); j++ {
		cPorts := []v1.ContainerPort{
			{
				Name:          pods.Container.ContainerPort[j].Name,
				ContainerPort: pods.Container.ContainerPort[j].Port,
			},
		}
		ports = append(ports, cPorts...)
	}

	for k := 0; k < len(pods.Container.VolumeMount); k++ {
		mounts := []v1.VolumeMount{
			{
				MountPath: pods.Container.VolumeMount[k].MountPath,
				Name:      pods.Container.VolumeMount[k].MountName,
			},
		}
		mountPath = append(mountPath, mounts...)
	}

	for l := 0; l < len(pods.Container.EnvVariable); l++ {
		envOut := []v1.EnvVar{
			{
				Name:  pods.Container.EnvVariable[l].Name,
				Value: pods.Container.EnvVariable[l].Value,
			},
		}
		env = append(env, envOut...)
	}

	for a := 0; a < len(pods.Container.EnvFromSecret); a++ {
		envOut := []v1.EnvVar{
			{
				Name: pods.Container.EnvFromSecret[a].Name,
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: pods.Container.EnvFromSecret[a].SecretName,
						},
						Key: pods.Container.EnvFromSecret[a].SecretKey,
					},
				},
			},
		}
		secretEnv = append(secretEnv, envOut...)
	}

	for b := 0; b < len(pods.Container.EnvFromConfigmap); b++ {
		envOut := []v1.EnvVar{
			{
				Name: pods.Container.EnvFromConfigmap[b].Name,
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: pods.Container.EnvFromConfigmap[b].ConfigmapName,
						},
						Key: pods.Container.EnvFromConfigmap[b].ConfigmapKey,
					},
				},
			},
		}
		configEnv = append(configEnv, envOut...)
	}

	// Container Cpu and Memory Resources

	if pods.Container.CpuLimits != 0 && pods.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
		}
	} else if pods.Container.CpuRequest != 0 && pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits != 0 {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
			},
		}
	} else if pods.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
		}
	} else if pods.Container.CpuRequest != 0 {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits != 0 && pods.Container.MemoryLimits == "" && pods.Container.CpuRequest != 0 && pods.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if pods.Container.CpuLimits == 0 && pods.Container.MemoryLimits != "" && pods.Container.CpuRequest == 0 && pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits != 0 && pods.Container.MemoryLimits == "" && pods.Container.CpuRequest != 0 && pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits == 0 && pods.Container.MemoryLimits != "" && pods.Container.CpuRequest != 0 && pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits != 0 && pods.Container.MemoryLimits != "" && pods.Container.CpuRequest == 0 && pods.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{

				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	} else if pods.Container.CpuLimits != 0 && pods.Container.MemoryLimits != "" && pods.Container.CpuRequest != 0 && pods.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(pods.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(pods.Container.MemoryRequest),
			},
		}
	}

	// Init Container

	if pods.InitContainer.Command != nil && pods.InitContainer.Args == nil {
		initContainers = []v1.Container{
			{
				Name:         pods.InitContainer.ContainerName,
				Image:        pods.InitContainer.Image,
				Command:      pods.InitContainer.Command,
				VolumeMounts: mountPath,
			},
		}
	} else if pods.InitContainer.Command == nil && pods.InitContainer.Args != nil {
		initContainers = []v1.Container{
			{
				Name:         pods.InitContainer.ContainerName,
				Image:        pods.InitContainer.Image,
				Args:         pods.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	} else {
		initContainers = []v1.Container{
			{
				Name:         pods.InitContainer.ContainerName,
				Image:        pods.InitContainer.Image,
				Command:      pods.InitContainer.Command,
				Args:         pods.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	}

	// Pod Security Context

	if pods.SecurityContext.RunAsUser != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &pods.SecurityContext.RunAsUser,
		}
	} else if pods.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &pods.SecurityContext.RunAsGroup,
		}
	} else if pods.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			FSGroup: &pods.SecurityContext.FsGroup,
		}
	} else if pods.SecurityContext.RunAsUser != 0 && pods.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &pods.SecurityContext.RunAsGroup,
			RunAsUser:  &pods.SecurityContext.RunAsUser,
		}
	} else if pods.SecurityContext.RunAsUser != 0 && pods.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &pods.SecurityContext.RunAsUser,
			FSGroup:   &pods.SecurityContext.FsGroup,
		}
	} else if pods.SecurityContext.RunAsGroup != 0 && pods.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &pods.SecurityContext.RunAsGroup,
			FSGroup:   &pods.SecurityContext.FsGroup,
		}
	} else {
		security = v1.PodSecurityContext{
			RunAsUser:  &pods.SecurityContext.RunAsUser,
			RunAsGroup: &pods.SecurityContext.RunAsGroup,
			FSGroup:    &pods.SecurityContext.FsGroup,
		}
	}

	// Container Security Context

	if pods.Container.SecurityContext.RunAsUser != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsUser: &pods.Container.SecurityContext.RunAsUser,
		}
	} else if pods.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &pods.Container.SecurityContext.RunAsGroup,
		}
	} else if pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &pods.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &pods.Container.SecurityContext.RunAsUser,
		}
	} else {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &pods.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &pods.Container.SecurityContext.RunAsUser,
			Privileged: &pods.Container.SecurityContext.Privileged,
		}
	}

	// Liveness Probe

	liveness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: pods.Container.LivenessProbe.Command,
			},
		},
		InitialDelaySeconds: pods.Container.LivenessProbe.InitialDelaySeconds,
		PeriodSeconds:       pods.Container.LivenessProbe.PeriodSeconds,
	}

	// Readiness Probe

	readiness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: pods.Container.ReadinessProbe.Command,
			},
		},
		InitialDelaySeconds: pods.Container.ReadinessProbe.InitialDelaySeconds,
		PeriodSeconds:       pods.Container.ReadinessProbe.PeriodSeconds,
	}

	//Containers

	if pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.Command != nil && pods.Container.Args == nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Command:      pods.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}

	} else if pods.Container.Command != nil && pods.Container.Args == nil && pods.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Command:      pods.Container.Command,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.Command != nil && pods.Container.Args == nil && pods.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Command:      pods.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}

	} else if pods.Container.Command == nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Args:         pods.Container.Args,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.Command == nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Args:         pods.Container.Args,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}
	} else if pods.Container.Command == nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         pods.Container.ContainerName,
				Image:        pods.Container.Image,
				Ports:        ports,
				Args:         pods.Container.Args,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          pods.Container.Tty,
			},
		}

	} else if pods.Container.SecurityContext.RunAsGroup != 0 || pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.ReadinessProbe.Command != nil && pods.Container.LivenessProbe.Command == nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command != nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         pods.Container.Command,
				Args:            pods.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if pods.Container.SecurityContext.RunAsGroup != 0 || pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.ReadinessProbe.Command != nil && pods.Container.LivenessProbe.Command == nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command != nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         pods.Container.Command,
				Args:            pods.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if pods.Container.SecurityContext.RunAsGroup != 0 || pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.ReadinessProbe.Command == nil && pods.Container.LivenessProbe.Command != nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command != nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         pods.Container.Command,
				Args:            pods.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if pods.Container.SecurityContext.RunAsGroup != 0 || pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.ReadinessProbe.Command == nil && pods.Container.LivenessProbe.Command != nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command != nil && pods.Container.Args == nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         pods.Container.Command,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if pods.Container.SecurityContext.RunAsGroup != 0 || pods.Container.SecurityContext.RunAsUser != 0 && pods.Container.ReadinessProbe.Command != nil && pods.Container.LivenessProbe.Command != nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command == nil && pods.Container.Args != nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             env,
				Args:            pods.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if pods.Container.ReadinessProbe.Command != nil && pods.Container.LivenessProbe.Command != nil && pods.Container.CpuLimits != 0 || pods.Container.MemoryLimits != "" || pods.Container.CpuRequest != 0 || pods.Container.MemoryRequest != "" && pods.Container.Command != nil && pods.Container.Args == nil && pods.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:           pods.Container.ContainerName,
				Image:          pods.Container.Image,
				Ports:          ports,
				Env:            env,
				Command:        pods.Container.Command,
				Resources:      resources,
				LivenessProbe:  liveness,
				ReadinessProbe: readiness,
				VolumeMounts:   mountPath,
				TTY:            pods.Container.Tty,
			},
		}
	} else {
		cnt = []v1.Container{
			{
				Name:            pods.Container.ContainerName,
				Image:           pods.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         pods.Container.Command,
				Args:            pods.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             pods.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	}

	if pods.VolumeSource == "HostPath" || pods.VolumeSource == "Hostpath" || pods.VolumeSource == "hostpath" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers: cnt,
				NodeName:   pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.HostPath.MountName,
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: pods.Storage.HostPath.Path,
								Type: (*v1.HostPathType)(&pods.Storage.HostPath.Type),
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "NFS" || pods.VolumeSource == "Nfs" || pods.VolumeSource == "nfs" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers: cnt,
				NodeName:   pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.NFS.MountName,
						VolumeSource: v1.VolumeSource{
							NFS: &v1.NFSVolumeSource{
								Server: pods.Storage.NFS.Server,
								Path:   pods.Storage.NFS.Path,
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers: cnt,
				NodeName:   pods.NodeName,
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "pvc" || pods.VolumeSource == "PVC" || pods.VolumeSource == "Pvc" && pods.InitContainer.ContainerName != "" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers:     cnt,
				InitContainers: initContainers,
				NodeName:       pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.PVC.MountName,
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: pods.Storage.PVC.ClaimName,
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "pvc" || pods.VolumeSource == "PVC" || pods.VolumeSource == "Pvc" && pods.InitContainer.ContainerName == "" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers: cnt,
				NodeName:   pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.PVC.MountName,
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: pods.Storage.PVC.ClaimName,
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "pvcconfig" || pods.VolumeSource == "PVCCONFIG" || pods.VolumeSource == "Pvcconfig" && pods.InitContainer.ContainerName != "" && pods.SecurityContext.RunAsUser != 0 {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				SecurityContext: &security,
				Containers:      cnt,
				InitContainers:  initContainers,
				NodeName:        pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.PVC.MountName,
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: pods.Storage.PVC.ClaimName,
							},
						},
					},
					{
						Name: pods.Storage.ConfigmapSourceOne.MountName,
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: pods.Storage.ConfigmapSourceOne.ConfigmapName,
								},
								Items: []v1.KeyToPath{
									{
										Key:  pods.Storage.ConfigmapSourceOne.Key,
										Path: pods.Storage.ConfigmapSourceOne.Path,
									},
								},
								DefaultMode: pods.Storage.ConfigmapSourceOne.Mode,
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if pods.VolumeSource == "pvcconfig" || pods.VolumeSource == "PVCCONFIG" || pods.VolumeSource == "Pvcconfig" && pods.InitContainer.ContainerName == "" {
		podCreate := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pods.Name,
				Namespace: pods.Namespace,
				Labels:    labels,
			},
			Spec: v1.PodSpec{
				Containers:     cnt,
				InitContainers: initContainers,
				NodeName:       pods.NodeName,
				Volumes: []v1.Volume{
					{
						Name: pods.Storage.PVC.MountName,
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: pods.Storage.PVC.ClaimName,
							},
						},
					},
					{
						Name: pods.Storage.ConfigmapSourceOne.MountName,
						VolumeSource: v1.VolumeSource{
							ConfigMap: &v1.ConfigMapVolumeSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: pods.Storage.ConfigmapSourceOne.ConfigmapName,
								},
								Items: []v1.KeyToPath{
									{
										Key:  pods.Storage.ConfigmapSourceOne.Key,
										Path: pods.Storage.ConfigmapSourceOne.Path,
									},
								},
								DefaultMode: pods.Storage.ConfigmapSourceOne.Mode,
							},
						},
					},
				},
			},
		}
		create, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), podCreate, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	}
}

func (u UserImpl) ListPods(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetPod(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeletePod(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// Deployments

func (u UserImpl) CreateDeployment(namespace string, dep Models.Deployment, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&dep)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(dep.ClusterId, w, r)
	var labels map[string]string
	var cnt []v1.Container
	var ports []v1.ContainerPort
	var env []v1.EnvVar
	var mountPath []v1.VolumeMount
	var resources v1.ResourceRequirements
	var initContainers []v1.Container
	var secretEnv []v1.EnvVar
	var configEnv []v1.EnvVar
	var security v1.PodSecurityContext
	var cntsecurity v1.SecurityContext
	var readiness *v1.Probe
	var liveness *v1.Probe

	labels = make(map[string]string)
	for i := 0; i < len(dep.Labels); i++ {
		labels[dep.Labels[i].FirstLabel] = dep.Labels[i].SecondLabel
	}

	for j := 0; j < len(dep.Container.ContainerPort); j++ {
		cPorts := []v1.ContainerPort{
			{
				Name:          dep.Container.ContainerPort[j].Name,
				ContainerPort: dep.Container.ContainerPort[j].Port,
			},
		}
		ports = append(ports, cPorts...)
	}

	for k := 0; k < len(dep.Container.VolumeMount); k++ {
		mounts := []v1.VolumeMount{
			{
				MountPath: dep.Container.VolumeMount[k].MountPath,
				Name:      dep.Container.VolumeMount[k].MountName,
			},
		}
		mountPath = append(mountPath, mounts...)
	}

	for l := 0; l < len(dep.Container.EnvVariable); l++ {
		envOut := []v1.EnvVar{
			{
				Name:  dep.Container.EnvVariable[l].Name,
				Value: dep.Container.EnvVariable[l].Value,
			},
		}
		env = append(env, envOut...)
	}

	for a := 0; a < len(dep.Container.EnvFromSecret); a++ {
		envOut := []v1.EnvVar{
			{
				Name: dep.Container.EnvFromSecret[a].Name,
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: dep.Container.EnvFromSecret[a].SecretName,
						},
						Key: dep.Container.EnvFromSecret[a].SecretKey,
					},
				},
			},
		}
		secretEnv = append(secretEnv, envOut...)
	}

	for b := 0; b < len(dep.Container.EnvFromConfigmap); b++ {
		envOut := []v1.EnvVar{
			{
				Name: dep.Container.EnvFromConfigmap[b].Name,
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: dep.Container.EnvFromConfigmap[b].ConfigmapName,
						},
						Key: dep.Container.EnvFromConfigmap[b].ConfigmapKey,
					},
				},
			},
		}
		configEnv = append(configEnv, envOut...)
	}

	// Container Cpu and Memory Resources

	if dep.Container.CpuLimits != 0 && dep.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
		}
	} else if dep.Container.CpuRequest != 0 && dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits != 0 {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
			},
		}
	} else if dep.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
		}
	} else if dep.Container.CpuRequest != 0 {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits != 0 && dep.Container.MemoryLimits == "" && dep.Container.CpuRequest != 0 && dep.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if dep.Container.CpuLimits == 0 && dep.Container.MemoryLimits != "" && dep.Container.CpuRequest == 0 && dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits != 0 && dep.Container.MemoryLimits == "" && dep.Container.CpuRequest != 0 && dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits == 0 && dep.Container.MemoryLimits != "" && dep.Container.CpuRequest != 0 && dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits != 0 && dep.Container.MemoryLimits != "" && dep.Container.CpuRequest == 0 && dep.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{

				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	} else if dep.Container.CpuLimits != 0 && dep.Container.MemoryLimits != "" && dep.Container.CpuRequest != 0 && dep.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(dep.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(dep.Container.MemoryRequest),
			},
		}
	}

	// Init Container

	if dep.InitContainer.Command != nil && dep.InitContainer.Args == nil {
		initContainers = []v1.Container{
			{
				Name:         dep.InitContainer.ContainerName,
				Image:        dep.InitContainer.Image,
				Command:      dep.InitContainer.Command,
				VolumeMounts: mountPath,
			},
		}
	} else if dep.InitContainer.Command == nil && dep.InitContainer.Args != nil {
		initContainers = []v1.Container{
			{
				Name:         dep.InitContainer.ContainerName,
				Image:        dep.InitContainer.Image,
				Args:         dep.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	} else {
		initContainers = []v1.Container{
			{
				Name:         dep.InitContainer.ContainerName,
				Image:        dep.InitContainer.Image,
				Command:      dep.InitContainer.Command,
				Args:         dep.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	}

	// Pod Security Context

	if dep.SecurityContext.RunAsUser != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &dep.SecurityContext.RunAsUser,
		}
	} else if dep.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &dep.SecurityContext.RunAsGroup,
		}
	} else if dep.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			FSGroup: &dep.SecurityContext.FsGroup,
		}
	} else if dep.SecurityContext.RunAsUser != 0 && dep.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &dep.SecurityContext.RunAsGroup,
			RunAsUser:  &dep.SecurityContext.RunAsUser,
		}
	} else if dep.SecurityContext.RunAsUser != 0 && dep.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &dep.SecurityContext.RunAsUser,
			FSGroup:   &dep.SecurityContext.FsGroup,
		}
	} else if dep.SecurityContext.RunAsGroup != 0 && dep.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &dep.SecurityContext.RunAsGroup,
			FSGroup:   &dep.SecurityContext.FsGroup,
		}
	} else {
		security = v1.PodSecurityContext{
			RunAsUser:  &dep.SecurityContext.RunAsUser,
			RunAsGroup: &dep.SecurityContext.RunAsGroup,
			FSGroup:    &dep.SecurityContext.FsGroup,
		}
	}

	// Container Security Context

	if dep.Container.SecurityContext.RunAsUser != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsUser: &dep.Container.SecurityContext.RunAsUser,
		}
	} else if dep.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &dep.Container.SecurityContext.RunAsGroup,
		}
	} else if dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &dep.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &dep.Container.SecurityContext.RunAsUser,
		}
	} else {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &dep.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &dep.Container.SecurityContext.RunAsUser,
			Privileged: &dep.Container.SecurityContext.Privileged,
		}
	}

	// Liveness Probe

	liveness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: dep.Container.LivenessProbe.Command,
			},
		},
		InitialDelaySeconds: dep.Container.LivenessProbe.InitialDelaySeconds,
		PeriodSeconds:       dep.Container.LivenessProbe.PeriodSeconds,
	}

	// Readiness Probe

	readiness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: dep.Container.ReadinessProbe.Command,
			},
		},
		InitialDelaySeconds: dep.Container.ReadinessProbe.InitialDelaySeconds,
		PeriodSeconds:       dep.Container.ReadinessProbe.PeriodSeconds,
	}

	//Containers

	if dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.Command != nil && dep.Container.Args == nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Command:      dep.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}

	} else if dep.Container.Command != nil && dep.Container.Args == nil && dep.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Command:      dep.Container.Command,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.Command != nil && dep.Container.Args == nil && dep.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Command:      dep.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}

	} else if dep.Container.Command == nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Args:         dep.Container.Args,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.Command == nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Args:         dep.Container.Args,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}
	} else if dep.Container.Command == nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         dep.Container.ContainerName,
				Image:        dep.Container.Image,
				Ports:        ports,
				Args:         dep.Container.Args,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          dep.Container.Tty,
			},
		}

	} else if dep.Container.SecurityContext.RunAsGroup != 0 || dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.ReadinessProbe.Command != nil && dep.Container.LivenessProbe.Command == nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command != nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         dep.Container.Command,
				Args:            dep.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if dep.Container.SecurityContext.RunAsGroup != 0 || dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.ReadinessProbe.Command != nil && dep.Container.LivenessProbe.Command == nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command != nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         dep.Container.Command,
				Args:            dep.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if dep.Container.SecurityContext.RunAsGroup != 0 || dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.ReadinessProbe.Command == nil && dep.Container.LivenessProbe.Command != nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command != nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         dep.Container.Command,
				Args:            dep.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if dep.Container.SecurityContext.RunAsGroup != 0 || dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.ReadinessProbe.Command == nil && dep.Container.LivenessProbe.Command != nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command != nil && dep.Container.Args == nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         dep.Container.Command,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if dep.Container.SecurityContext.RunAsGroup != 0 || dep.Container.SecurityContext.RunAsUser != 0 && dep.Container.ReadinessProbe.Command != nil && dep.Container.LivenessProbe.Command != nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command == nil && dep.Container.Args != nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             env,
				Args:            dep.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if dep.Container.ReadinessProbe.Command != nil && dep.Container.LivenessProbe.Command != nil && dep.Container.CpuLimits != 0 || dep.Container.MemoryLimits != "" || dep.Container.CpuRequest != 0 || dep.Container.MemoryRequest != "" && dep.Container.Command != nil && dep.Container.Args == nil && dep.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:           dep.Container.ContainerName,
				Image:          dep.Container.Image,
				Ports:          ports,
				Env:            env,
				Command:        dep.Container.Command,
				Resources:      resources,
				LivenessProbe:  liveness,
				ReadinessProbe: readiness,
				VolumeMounts:   mountPath,
				TTY:            dep.Container.Tty,
			},
		}
	} else {
		cnt = []v1.Container{
			{
				Name:            dep.Container.ContainerName,
				Image:           dep.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         dep.Container.Command,
				Args:            dep.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             dep.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	}

	if dep.VolumeSource == "PVC" || dep.VolumeSource == "Pvc" || dep.VolumeSource == "pvc" && dep.InitContainer.Image != "" && dep.SecurityContext.RunAsUser != 0 || dep.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						InitContainers:  initContainers,
						SecurityContext: &security,
						NodeName:        dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "PVC" || dep.VolumeSource == "Pvc" || dep.VolumeSource == "pvc" && dep.InitContainer.Image == "" && dep.SecurityContext.RunAsUser != 0 || dep.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						SecurityContext: &security,
						NodeName:        dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "PVC" || dep.VolumeSource == "Pvc" || dep.VolumeSource == "pvc" {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "HostPath" || dep.VolumeSource == "Hostpath" || dep.VolumeSource == "hostpath" {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.HostPath.MountName,
								VolumeSource: v1.VolumeSource{
									HostPath: &v1.HostPathVolumeSource{
										Path: dep.Storage.HostPath.Path,
										Type: (*v1.HostPathType)(&dep.Storage.HostPath.Type),
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "NFS" || dep.VolumeSource == "Nfs" || dep.VolumeSource == "nfs" {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.NFS.MountName,
								VolumeSource: v1.VolumeSource{
									NFS: &v1.NFSVolumeSource{
										Server: dep.Storage.NFS.Server,
										Path:   dep.Storage.NFS.Path,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "Secret" || dep.VolumeSource == "secret" {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.SecretSource.MountName,
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName: dep.Storage.SecretSource.SecretName,
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.SecretSource.Key,
												Path: dep.Storage.SecretSource.Path,
											},
										},
										DefaultMode: dep.Storage.SecretSource.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "Configmap" || dep.VolumeSource == "ConfigMap" || dep.VolumeSource == "configmap" {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceOne.Key,
												Path: dep.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "pvcconfig" || dep.VolumeSource == "PVCCONFIG" || dep.VolumeSource == "Pvcconfig" && dep.InitContainer.ContainerName != "" && dep.SecurityContext.RunAsUser != 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						SecurityContext: &security,
						Containers:      cnt,
						InitContainers:  initContainers,
						NodeName:        dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: dep.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceOne.Key,
												Path: dep.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
						},
					},
				},
			},
		}

		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "pvcconfig" || dep.VolumeSource == "PVCCONFIG" || dep.VolumeSource == "Pvcconfig" && dep.InitContainer.ContainerName == "" && dep.SecurityContext.RunAsUser == 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: dep.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceOne.Key,
												Path: dep.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
						},
					},
				},
			},
		}

		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "mysql" || dep.VolumeSource == "Mysql" || dep.VolumeSource == "MySqlCluster" && dep.InitContainer.ContainerName != "" && dep.SecurityContext.RunAsUser != 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						InitContainers:  initContainers,
						SecurityContext: &security,
						NodeName:        dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceOne.Key,
												Path: dep.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: dep.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceTwo.Key,
												Path: dep.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if dep.VolumeSource == "mysql" || dep.VolumeSource == "Mysql" || dep.VolumeSource == "MySqlCluster" && dep.InitContainer.ContainerName == "" && dep.SecurityContext.RunAsUser == 0 {
		deploy := &v2.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      dep.Name,
				Namespace: dep.Namespace,
				Labels:    labels,
			},
			Spec: v2.DeploymentSpec{
				Replicas: &dep.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   dep.NodeName,
						Volumes: []v1.Volume{
							{
								Name: dep.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceOne.Key,
												Path: dep.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: dep.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: dep.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: dep.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: dep.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  dep.Storage.ConfigmapSourceTwo.Key,
												Path: dep.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: dep.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	}
}

func (u UserImpl) ListDeployment(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetDeployment(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeleteDeployment(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// Statefulsets

func (u UserImpl) CreateStatefulSet(namespace string, sts Models.StatefulSet, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&sts)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(sts.ClusterId, w, r)
	var labels map[string]string
	var cnt []v1.Container
	var ports []v1.ContainerPort
	var env []v1.EnvVar
	var resources v1.ResourceRequirements
	var mountPath []v1.VolumeMount
	var initContainers []v1.Container
	var secretEnv []v1.EnvVar
	var configEnv []v1.EnvVar
	var security v1.PodSecurityContext
	var cntsecurity v1.SecurityContext
	var readiness *v1.Probe
	var liveness *v1.Probe

	labels = make(map[string]string)
	for i := 0; i < len(sts.Labels); i++ {
		labels[sts.Labels[i].FirstLabel] = sts.Labels[i].SecondLabel
	}

	for j := 0; j < len(sts.Container.ContainerPort); j++ {
		cPorts := []v1.ContainerPort{
			{
				Name:          sts.Container.ContainerPort[j].Name,
				ContainerPort: sts.Container.ContainerPort[j].Port,
			},
		}
		ports = append(ports, cPorts...)
	}

	for k := 0; k < len(sts.Container.VolumeMount); k++ {
		mounts := []v1.VolumeMount{
			{
				MountPath: sts.Container.VolumeMount[k].MountPath,
				Name:      sts.Container.VolumeMount[k].MountName,
				SubPath:   sts.Container.VolumeMount[k].SubPath,
			},
		}
		mountPath = append(mountPath, mounts...)
	}

	for l := 0; l < len(sts.Container.EnvVariable); l++ {
		envOut := []v1.EnvVar{
			{
				Name:  sts.Container.EnvVariable[l].Name,
				Value: sts.Container.EnvVariable[l].Value,
			},
		}
		env = append(env, envOut...)
	}

	for a := 0; a < len(sts.Container.EnvFromSecret); a++ {
		envOut := []v1.EnvVar{
			{
				Name: sts.Container.EnvFromSecret[a].Name,
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: sts.Container.EnvFromSecret[a].SecretName,
						},
						Key: sts.Container.EnvFromSecret[a].SecretKey,
					},
				},
			},
		}
		secretEnv = append(secretEnv, envOut...)
	}

	for b := 0; b < len(sts.Container.EnvFromConfigmap); b++ {
		envOut := []v1.EnvVar{
			{
				Name: sts.Container.EnvFromConfigmap[b].Name,
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: sts.Container.EnvFromConfigmap[b].ConfigmapName,
						},
						Key: sts.Container.EnvFromConfigmap[b].ConfigmapKey,
					},
				},
			},
		}
		configEnv = append(configEnv, envOut...)
	}

	// Container Cpu and Memory Resources

	if sts.Container.CpuLimits != 0 && sts.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
		}
	} else if sts.Container.CpuRequest != 0 && sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits != 0 {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
			},
		}
	} else if sts.Container.MemoryLimits != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
		}
	} else if sts.Container.CpuRequest != 0 {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits != 0 && sts.Container.MemoryLimits == "" && sts.Container.CpuRequest != 0 && sts.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else if sts.Container.CpuLimits == 0 && sts.Container.MemoryLimits != "" && sts.Container.CpuRequest == 0 && sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits != 0 && sts.Container.MemoryLimits == "" && sts.Container.CpuRequest != 0 && sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits == 0 && sts.Container.MemoryLimits != "" && sts.Container.CpuRequest != 0 && sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits != 0 && sts.Container.MemoryLimits != "" && sts.Container.CpuRequest == 0 && sts.Container.MemoryRequest != "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{

				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	} else if sts.Container.CpuLimits != 0 && sts.Container.MemoryLimits != "" && sts.Container.CpuRequest != 0 && sts.Container.MemoryRequest == "" {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU: *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
			},
		}
	} else {
		resources = v1.ResourceRequirements{
			Limits: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuLimits, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryLimits),
			},
			Requests: map[v1.ResourceName]resource.Quantity{
				v1.ResourceCPU:    *resource.NewMilliQuantity(sts.Container.CpuRequest, resource.DecimalSI),
				v1.ResourceMemory: resource.MustParse(sts.Container.MemoryRequest),
			},
		}
	}

	// Init Container

	if sts.InitContainer.Command != nil && sts.InitContainer.Args == nil {
		initContainers = []v1.Container{
			{
				Name:         sts.InitContainer.ContainerName,
				Image:        sts.InitContainer.Image,
				Command:      sts.InitContainer.Command,
				VolumeMounts: mountPath,
			},
		}
	} else if sts.InitContainer.Command == nil && sts.InitContainer.Args != nil {
		initContainers = []v1.Container{
			{
				Name:         sts.InitContainer.ContainerName,
				Image:        sts.InitContainer.Image,
				Args:         sts.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	} else {
		initContainers = []v1.Container{
			{
				Name:         sts.InitContainer.ContainerName,
				Image:        sts.InitContainer.Image,
				Command:      sts.InitContainer.Command,
				Args:         sts.InitContainer.Args,
				VolumeMounts: mountPath,
			},
		}
	}

	// Pod Security Context

	if sts.SecurityContext.RunAsUser != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &sts.SecurityContext.RunAsUser,
		}
	} else if sts.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &sts.SecurityContext.RunAsGroup,
		}
	} else if sts.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			FSGroup: &sts.SecurityContext.FsGroup,
		}
	} else if sts.SecurityContext.RunAsUser != 0 && sts.SecurityContext.RunAsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsGroup: &sts.SecurityContext.RunAsGroup,
			RunAsUser:  &sts.SecurityContext.RunAsUser,
		}
	} else if sts.SecurityContext.RunAsUser != 0 && sts.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &sts.SecurityContext.RunAsUser,
			FSGroup:   &sts.SecurityContext.FsGroup,
		}
	} else if sts.SecurityContext.RunAsGroup != 0 && sts.SecurityContext.FsGroup != 0 {
		security = v1.PodSecurityContext{
			RunAsUser: &sts.SecurityContext.RunAsGroup,
			FSGroup:   &sts.SecurityContext.FsGroup,
		}
	} else {
		security = v1.PodSecurityContext{
			RunAsUser:  &sts.SecurityContext.RunAsUser,
			RunAsGroup: &sts.SecurityContext.RunAsGroup,
			FSGroup:    &sts.SecurityContext.FsGroup,
		}
	}

	// Container Security Context

	if sts.Container.SecurityContext.RunAsUser != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsUser: &sts.Container.SecurityContext.RunAsUser,
		}
	} else if sts.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &sts.Container.SecurityContext.RunAsGroup,
		}
	} else if sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.SecurityContext.RunAsGroup != 0 {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &sts.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &sts.Container.SecurityContext.RunAsUser,
		}
	} else {
		cntsecurity = v1.SecurityContext{
			RunAsGroup: &sts.Container.SecurityContext.RunAsGroup,
			RunAsUser:  &sts.Container.SecurityContext.RunAsUser,
			Privileged: &sts.Container.SecurityContext.Privileged,
		}
	}

	// Liveness Probe

	liveness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: sts.Container.LivenessProbe.Command,
			},
		},
		InitialDelaySeconds: sts.Container.LivenessProbe.InitialDelaySeconds,
		PeriodSeconds:       sts.Container.LivenessProbe.PeriodSeconds,
	}

	// Readiness Probe

	readiness = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &v1.ExecAction{
				Command: sts.Container.ReadinessProbe.Command,
			},
		},
		InitialDelaySeconds: sts.Container.ReadinessProbe.InitialDelaySeconds,
		PeriodSeconds:       sts.Container.ReadinessProbe.PeriodSeconds,
	}

	//Containers

	if sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.Command != nil && sts.Container.Args == nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Command:      sts.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}

	} else if sts.Container.Command != nil && sts.Container.Args == nil && sts.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Command:      sts.Container.Command,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.Command != nil && sts.Container.Args == nil && sts.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Command:      sts.Container.Command,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}

	} else if sts.Container.Command == nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Args:         sts.Container.Args,
				Env:          env,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.Command == nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Args:         sts.Container.Args,
				Env:          secretEnv,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}
	} else if sts.Container.Command == nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvFromConfigmap" {
		cnt = []v1.Container{
			{
				Name:         sts.Container.ContainerName,
				Image:        sts.Container.Image,
				Ports:        ports,
				Args:         sts.Container.Args,
				Env:          configEnv,
				VolumeMounts: mountPath,
				TTY:          sts.Container.Tty,
			},
		}

	} else if sts.Container.SecurityContext.RunAsGroup != 0 || sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.ReadinessProbe.Command != nil && sts.Container.LivenessProbe.Command == nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command != nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         sts.Container.Command,
				Args:            sts.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if sts.Container.SecurityContext.RunAsGroup != 0 || sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.ReadinessProbe.Command != nil && sts.Container.LivenessProbe.Command == nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command != nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         sts.Container.Command,
				Args:            sts.Container.Args,
				Resources:       resources,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if sts.Container.SecurityContext.RunAsGroup != 0 || sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.ReadinessProbe.Command == nil && sts.Container.LivenessProbe.Command != nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command != nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvFromSecret" {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             secretEnv,
				Command:         sts.Container.Command,
				Args:            sts.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if sts.Container.SecurityContext.RunAsGroup != 0 || sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.ReadinessProbe.Command == nil && sts.Container.LivenessProbe.Command != nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command != nil && sts.Container.Args == nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         sts.Container.Command,
				Resources:       resources,
				LivenessProbe:   liveness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if sts.Container.SecurityContext.RunAsGroup != 0 || sts.Container.SecurityContext.RunAsUser != 0 && sts.Container.ReadinessProbe.Command != nil && sts.Container.LivenessProbe.Command != nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command == nil && sts.Container.Args != nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             env,
				Args:            sts.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	} else if sts.Container.ReadinessProbe.Command != nil && sts.Container.LivenessProbe.Command != nil && sts.Container.CpuLimits != 0 || sts.Container.MemoryLimits != "" || sts.Container.CpuRequest != 0 || sts.Container.MemoryRequest != "" && sts.Container.Command != nil && sts.Container.Args == nil && sts.Container.EnvType == "EnvVar" {
		cnt = []v1.Container{
			{
				Name:           sts.Container.ContainerName,
				Image:          sts.Container.Image,
				Ports:          ports,
				Env:            env,
				Command:        sts.Container.Command,
				Resources:      resources,
				LivenessProbe:  liveness,
				ReadinessProbe: readiness,
				VolumeMounts:   mountPath,
				TTY:            sts.Container.Tty,
			},
		}
	} else {
		cnt = []v1.Container{
			{
				Name:            sts.Container.ContainerName,
				Image:           sts.Container.Image,
				Ports:           ports,
				Env:             env,
				Command:         sts.Container.Command,
				Args:            sts.Container.Args,
				Resources:       resources,
				LivenessProbe:   liveness,
				ReadinessProbe:  readiness,
				VolumeMounts:    mountPath,
				TTY:             sts.Container.Tty,
				SecurityContext: &cntsecurity,
			},
		}
	}

	if sts.VolumeSource == "PVC" || sts.VolumeSource == "Pvc" || sts.VolumeSource == "pvc" && sts.InitContainer.Image != "" && sts.SecurityContext.RunAsUser != 0 || sts.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						InitContainers:  initContainers,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "PVC" || sts.VolumeSource == "Pvc" || sts.VolumeSource == "pvc" && sts.InitContainer.Image == "" && sts.SecurityContext.RunAsUser != 0 || sts.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "PVC" || sts.VolumeSource == "Pvc" || sts.VolumeSource == "pvc" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "HostPath" || sts.VolumeSource == "Hostpath" || sts.VolumeSource == "hostpath" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.HostPath.MountName,
								VolumeSource: v1.VolumeSource{
									HostPath: &v1.HostPathVolumeSource{
										Path: sts.Storage.HostPath.Path,
										Type: (*v1.HostPathType)(&sts.Storage.HostPath.Type),
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "NFS" || sts.VolumeSource == "Nfs" || sts.VolumeSource == "nfs" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.NFS.MountName,
								VolumeSource: v1.VolumeSource{
									NFS: &v1.NFSVolumeSource{
										Server: sts.Storage.NFS.Server,
										Path:   sts.Storage.NFS.Path,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "Secret" || sts.VolumeSource == "secret" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.SecretSource.MountName,
								VolumeSource: v1.VolumeSource{
									Secret: &v1.SecretVolumeSource{
										SecretName: sts.Storage.SecretSource.SecretName,
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.SecretSource.Key,
												Path: sts.Storage.SecretSource.Path,
											},
										},
										DefaultMode: sts.Storage.SecretSource.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "Configmap" || sts.VolumeSource == "ConfigMap" || sts.VolumeSource == "configmap" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "VolumeClaimTemplate" || sts.VolumeSource == "Volumeclaimtemplate" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				VolumeClaimTemplates: []v1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: sts.Storage.VolumeClaimTemplate.Name,
						},
						Spec: v1.PersistentVolumeClaimSpec{
							AccessModes: []v1.PersistentVolumeAccessMode{
								v1.PersistentVolumeAccessMode(sts.Storage.VolumeClaimTemplate.AccessMode),
							},
							StorageClassName: &sts.Storage.VolumeClaimTemplate.StorageClassName,
							Resources: v1.ResourceRequirements{
								Requests: map[v1.ResourceName]resource.Quantity{
									v1.ResourceStorage: resource.MustParse(sts.Storage.VolumeClaimTemplate.Capacity),
								},
							},
						},
					},
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "pvcConfig" || sts.VolumeSource == "pvcconfig" || sts.VolumeSource == "PvcConfig" && sts.InitContainer.Image != "" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:     cnt,
						InitContainers: initContainers,
						NodeName:       sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "pvcConfig" || sts.VolumeSource == "pvcconfig" || sts.VolumeSource == "PvcConfig" && sts.InitContainer.Image != "" && sts.SecurityContext.RunAsGroup != 0 || sts.SecurityContext.RunAsUser != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						InitContainers:  initContainers,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "pvcConfig" || sts.VolumeSource == "pvcconfig" || sts.VolumeSource == "PvcConfig" && sts.SecurityContext.RunAsGroup != 0 || sts.SecurityContext.RunAsUser != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "mysql" || sts.VolumeSource == "Mysql" || sts.VolumeSource == "MySqlCluster" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers: cnt,
						NodeName:   sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: sts.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceTwo.Key,
												Path: sts.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "mysql" || sts.VolumeSource == "Mysql" || sts.VolumeSource == "MySqlCluster" && sts.InitContainer.Image != "" {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:     cnt,
						InitContainers: initContainers,
						NodeName:       sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: sts.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceTwo.Key,
												Path: sts.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "mysql" || sts.VolumeSource == "Mysql" || sts.VolumeSource == "MySqlCluster" && sts.InitContainer.Image != "" && sts.SecurityContext.RunAsUser != 0 || sts.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						InitContainers:  initContainers,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: sts.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceTwo.Key,
												Path: sts.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	} else if sts.VolumeSource == "mysql" || sts.VolumeSource == "Mysql" || sts.VolumeSource == "MySqlCluster" && sts.SecurityContext.RunAsUser != 0 || sts.SecurityContext.RunAsGroup != 0 {
		deploy := &v2.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sts.Name,
				Namespace: sts.Namespace,
				Labels:    labels,
			},
			Spec: v2.StatefulSetSpec{
				Replicas:    &sts.Replicas,
				ServiceName: sts.ServiceName,
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: v1.PodSpec{
						Containers:      cnt,
						SecurityContext: &security,
						NodeName:        sts.NodeName,
						Volumes: []v1.Volume{
							{
								Name: sts.Storage.ConfigmapSourceOne.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceOne.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceOne.Key,
												Path: sts.Storage.ConfigmapSourceOne.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceOne.Mode,
									},
								},
							},
							{
								Name: sts.Storage.PVC.MountName,
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: sts.Storage.PVC.ClaimName,
									},
								},
							},
							{
								Name: sts.Storage.ConfigmapSourceTwo.MountName,
								VolumeSource: v1.VolumeSource{
									ConfigMap: &v1.ConfigMapVolumeSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: sts.Storage.ConfigmapSourceTwo.ConfigmapName,
										},
										Items: []v1.KeyToPath{
											{
												Key:  sts.Storage.ConfigmapSourceTwo.Key,
												Path: sts.Storage.ConfigmapSourceTwo.Path,
											},
										},
										DefaultMode: sts.Storage.ConfigmapSourceTwo.Mode,
									},
								},
							},
						},
					},
				},
			},
		}
		create, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			logrus.Error("ResponseCode:", 2018, "Results:", err)
			response.RespondJSON(w, 2018, err)
		} else {
			response.RespondJSON(w, 1000, create)
		}
	}

}

func (u UserImpl) ListStatefulSets(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetStatefulSets(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeleteStatefulSets(w http.ResponseWriter, r *http.Request, cluster_id string, namespace string, name string) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

//Cluster Roles

func (u UserImpl) CreateClusterRole(cr Models.ClusterRole, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(cr.ClusterId, w, r)
	var labels map[string]string
	var clusterRole []av1.PolicyRule
	labels = make(map[string]string)
	for i := 0; i < len(cr.Labels); i++ {
		labels[cr.Labels[i].FirstLabel] = cr.Labels[i].SecondLabel
	}

	for j := 0; j < len(cr.RoleRules); j++ {
		cluster := []av1.PolicyRule{
			{
				APIGroups: cr.RoleRules[j].ApiGroup,
				Resources: cr.RoleRules[j].Resources,
				Verbs:     cr.RoleRules[j].Verbs,
			},
		}
		clusterRole = append(clusterRole, cluster...)
	}

	roleCreate := &av1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.Name,
			Labels: labels,
		},
		Rules: clusterRole,
	}

	create, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), roleCreate, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, create)
	}
}

func (u UserImpl) ListClusterRoles(cluster_id string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) GetClusterRole(cluster_id string, name string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) DeleteClusterRole(cluster_id string, name string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

// ClusterRole Binding

func (u UserImpl) GetClusterRoleBinding(cluster_id string, name string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	get, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, get)
	}
}

func (u UserImpl) ListClusterRoleBindings(cluster_id string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GettingCredentialsBasedOnClusterId(cluster_id, w, r)
	list, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, list)
	}
}

func (u UserImpl) DeleteClusterRoleBinding(cluster_id string, name string, w http.ResponseWriter, r *http.Request) {
	err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, nil)
	}
}

func (u UserImpl) CreateClusterRoleBinding(crb Models.ClusterRoleBinding, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&crb)
	if err != nil {
		logrus.Error("ResponseCode:", 2000, "Results:", err)
		response.RespondJSON(w, 2000, err)
	}
	GettingCredentialsBasedOnClusterId(crb.ClusterId, w, r)
	var labels map[string]string
	var subject []av1.Subject
	labels = make(map[string]string)
	for i := 0; i < len(crb.Labels); i++ {
		labels[crb.Labels[i].FirstLabel] = crb.Labels[i].SecondLabel
	}

	for j := 0; j < len(crb.Subject); j++ {
		sub := []av1.Subject{
			{
				Name:     crb.Subject[j].Name,
				Kind:     crb.Subject[j].Kind,
				APIGroup: crb.Subject[j].ApiGroup,
			},
		}
		subject = append(subject, sub...)
	}

	createCRB := &av1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   crb.Name,
			Labels: labels,
		},
		Subjects: subject,
		RoleRef: av1.RoleRef{
			APIGroup: crb.RoleRef.ApiGroup,
			Kind:     crb.RoleRef.Kind,
			Name:     crb.RoleRef.Name,
		},
	}
	create, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), createCRB, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("ResponseCode:", 2018, "Results:", err)
		response.RespondJSON(w, 2018, err)
	} else {
		response.RespondJSON(w, 1000, create)
	}
}
