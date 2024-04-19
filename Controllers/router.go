package Controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() {
	r := mux.NewRouter()

	//Checking Client Cluster Credentials
	r.HandleFunc("/api/kubernetes/checkCredentials", GetKubernetesCredentials).Methods("POST")

	//Cluster Nodes Details

	r.HandleFunc("/api/kubernetes/listNodes/{cluster_id}", ListClusterNodes).Methods("GET")
	r.HandleFunc("/api/kubernetes/getNode/{cluster_id}/{name}", GetClusterNodeDetail).Methods("GET")

	//Namespace

	r.HandleFunc("/api/kubernetes/createNamespace", CreateNamespace).Methods("POST")
	r.HandleFunc("/api/kubernetes/getNamespace/{cluster_id}/{name}", GetNamespace).Methods("GET")
	r.HandleFunc("/api/kubernetes/listNamespace/{cluster_id}", ListNamespace).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteNamespace/{cluster_id}/{name}", DeleteNamespace).Methods("DELETE")
	// r.HandleFunc("/api/kubernetes/updateNamespace", UpdateNamespace).Methods("PUT")

	//Service Account Routing

	r.HandleFunc("/api/kubernetes/createServiceAccount", CreateServiceAccount).Methods("POST")
	r.HandleFunc("/api/kubernetes/getServiceAccount/{cluster_id}/{namespace}/{name}", GetServiceAccount).Methods("GET")
	r.HandleFunc("/api/kubernetes/listServiceAccount/{cluster_id}/{namespace}", ListServiceAccount).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteServiceAccount/{cluster_id}/{namespace}/{name}", DeleteServiceAccountByName).Methods("DELETE")
	// r.HandleFunc("/api/kubernetes/updateServiceAccount/{namespace}", UpdateServiceAccountByName).Methods("PUT")

	// //ConfigMap
	r.HandleFunc("/api/kubernetes/createConfigmap/{namespace}", CreateConfigmapInCluster).Methods("POST")
	r.HandleFunc("/api/kubernetes/listConfigmap/{cluster_id}/{namespace}", ListConfigmapInCluster).Methods("GET")
	r.HandleFunc("/api/kubernetes/getConfigmapDetails/{cluster_id}/{namespace}/{name}", GetConfigmapDetails).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteConfigmap/{cluster_id}/{namespace}/{name}", DeleteConfigmap).Methods("DELETE")

	// StorageClass
	r.HandleFunc("/api/kubernetes/createStorageClass", CreateStorageClass).Methods("POST")
	r.HandleFunc("/api/kubernetes/listStorageClass/{cluster_id}", ListStorageClass).Methods("GET")
	r.HandleFunc("/api/kubernetes/getStorageClassDetails/{cluster_id}/{name}", GetDetailsOfStorageClass).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteStorageClass/{cluster_id}/{name}", DeleteStorageClass).Methods("DELETE")

	// Persistent Volume
	r.HandleFunc("/api/kubernetes/createPV", CreatePersistentVolume).Methods("POST")
	r.HandleFunc("/api/kubernetes/listPV/{cluster_id}", ListPersistentVolume).Methods("GET")
	r.HandleFunc("/api/kubernetes/getPV/{cluster_id}/{name}", GetDetailsOfAPersistentVolume).Methods("GET")
	r.HandleFunc("/api/kubernetes/deletePV/{cluster_id}/{name}", DeletePersistentVolume).Methods("DELETE")

	// Persistent Volume Claim
	r.HandleFunc("/api/kubernetes/createPVC/{namespace}", CreatePersistentVolumeClaim).Methods("POST")
	r.HandleFunc("/api/kubernetes/listPVC/{cluster_id}/{namespace}", ListPersistentVolumeClaim).Methods("GET")
	r.HandleFunc("/api/kubernetes/getPVC/{cluster_id}/{namespace}/{name}", GetPersistentVolumeClaimByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deletePVC/{cluster_id}/{namespace}/{name}", DeletePersistentVolumeClaimByName).Methods("DELETE")

	//Service
	r.HandleFunc("/api/kubernetes/createService/{namespace}", CreateService).Methods("POST")
	r.HandleFunc("/api/kubernetes/listServices/{cluster_id}/{namespace}", ListServiceInNamespace).Methods("GET")
	r.HandleFunc("/api/kubernetes/getService/{cluster_id}/{namespace}/{name}", GetServiceByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteService/{cluster_id}/{namespace}/{name}", DeleteServiceByName).Methods("DELETE")

	//Secret
	r.HandleFunc("/api/kubernetes/createSecret/{namespace}", CreateSecret).Methods("POST")
	r.HandleFunc("/api/kubernetes/getSecret/{cluster_id}/{namespace}/{name}", GetSecretByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/listSecrets/{cluster_id}/{namespace}", ListSecretsDetail).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteSecret/{cluster_id}/{namespace}/{name}", DeleteSecretByName).Methods("DELETE")

	//Pods
	r.HandleFunc("/api/kubernetes/createPod/{namespace}", CreatePod).Methods("POST")
	r.HandleFunc("/api/kubernetes/listPods/{cluster_id}/{namespace}", ListPodsInNamespace).Methods("GET")
	r.HandleFunc("/api/kubernetes/getPod/{cluster_id}/{namespace}/{name}", GetPodDetailsByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deletePod/{cluster_id}/{namespace}/{name}", DeletePodByName).Methods("DELETE")

	// //Deployment
	r.HandleFunc("/api/kubernetes/createDeployment/{namespace}", CreateDeployment).Methods("POST")
	r.HandleFunc("/api/kubernetes/listDeployments/{cluster_id}/{namespace}", ListDeploymentsDetail).Methods("GET")
	r.HandleFunc("/api/kubernetes/getDeployment/{cluster_id}/{namespace}/{name}", GetDeploymentDetailByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteDeployment/{cluster_id}/{namespace}/{name}", DeleteDeployment).Methods("DELETE")

	// //StatefulSet
	r.HandleFunc("/api/kubernetes/createStatefulset/{namespace}", CreateStatefulSet).Methods("POST")
	r.HandleFunc("/api/kubernetes/listStatefulset/{cluster_id}/{namespace}", ListStatefulSets).Methods("GET")
	r.HandleFunc("/api/kubernetes/getStatefulset/{cluster_id}/{namespace}/{name}", GetStatefulSetByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteStatefulset/{cluster_id}/{namespace}/{name}", DeleteStatefulSetByName).Methods("DELETE")

	//Cluster Role

	r.HandleFunc("/api/kubernetes/createClusterRole", CreateClusterRole).Methods("POST")
	r.HandleFunc("/api/kubernetes/listClusterRoles/{cluster_id}", ListClusterRoles).Methods("GET")
	r.HandleFunc("/api/kubernetes/getClusterRole/{cluster_id}/{name}", GetClusterRoleDetailsByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteClusterRole/{cluster_id}/{name}", DeleteClusterRolesByName).Methods("DELETE")

	//Cluster Role Bindings
	r.HandleFunc("/api/kubernetes/createClusterRoleBinding", CreateClusterRoleBinding).Methods("POST")
	r.HandleFunc("/api/kubernetes/listClusterRoleBindings/{cluster_id}", ListClusterRoleBindings).Methods("GET")
	r.HandleFunc("/api/kubernetes/getClusterRoleBinding/{cluster_id}/{name}", GetClusterRoleBindingDetailsByName).Methods("GET")
	r.HandleFunc("/api/kubernetes/deleteClusterRoleBinding/{cluster_id}/{name}", DeleteClusterRoleBindingByName).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":26443", r))
}
