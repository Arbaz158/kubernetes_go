package Controllers

import (
	"KubernetesGo/Dao"
	"KubernetesGo/Models"
	"net/http"

	"github.com/gorilla/mux"
)

var dao = Dao.UserImpl{}

func GetKubernetesCredentials(w http.ResponseWriter, r *http.Request) {
	var conn Models.Credentials
	Dao.CheckRemoteConnection(conn, w, r)
}

func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	var ns Models.Namespace
	dao.CreateNS(ns, w, r)
}

//Getting details about cluster

func ListClusterNodes(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListClusterNodes(w, r, id["cluster_id"])
}

func GetClusterNodeDetail(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	n := mux.Vars(r)
	dao.GetClusterNodes(w, r, id["cluster_id"], n["name"])
}

//Getting details about namespace

func GetNamespace(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetNs(w, r, id["cluster_id"], name["name"])
}

//Listing Namespace

func ListNamespace(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListNS(w, r, id["cluster_id"])
}

// //Deleting Namespace

func DeleteNamespace(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteNS(w, r, id["cluster_id"], params["name"])
}

// // Updating Namespace

// func UpdateNamespace(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var ns Models.Namespace
// 	err := json.NewDecoder(r.Body).Decode(&ns)
// 	if err != nil {
// 		fmt.Println("This error occurs inside Getting data from json body in UpdateNamespace function", err)
// 	}
// 	update, err := dao.UpdateNS(ns)
// 	if err != nil {
// 		fmt.Println("This error occurs inside UpdateNamespace function", err)
// 	}
// 	json.NewEncoder(w).Encode(update)
// }

func CreateServiceAccount(w http.ResponseWriter, r *http.Request) {
	var sa Models.ServiceAccount
	dao.CreateSA(w, r, sa)
}

func GetServiceAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ns := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetSA(w, r, id["cluster_id"], ns["namespace"], params["name"])
}

func ListServiceAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	params := mux.Vars(r)
	dao.ListSA(w, r, id["cluster_id"], params["namespace"])
}

func DeleteServiceAccountByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ns := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteSA(w, r, id["cluster_id"], ns["namespace"], params["name"])

}

// func UpdateServiceAccountByName(w http.ResponseWriter, r *http.Request)  {
// 	w.Header().Set("Content-Type", "application/json")
// 	var saupdate Models.ServiceAccount
// 	err := json.NewDecoder(r.Body).Decode(&saupdate)
// 	if err != nil {
// 		fmt.Println("This error occurs inside UpdateServiceAccountByName function", err)
// 	}
// 	params := mux.Vars(r)
// 	sa, err := dao.UpdateSA(saupdate, params["namespace"])
// 	if err != nil {
// 		fmt.Println("This is dao.Update error inside controller package", err)
// 	}
// 	json.NewEncoder(w).Encode(sa)
// }

//Configmaps

func CreateConfigmapInCluster(w http.ResponseWriter, r *http.Request) {
	var cfg Models.Configmap
	prj := mux.Vars(r)
	dao.CreateConfigmap(w, r, prj["namespace"], cfg)
}

func ListConfigmapInCluster(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListConfigmap(w, r, id["cluster_id"], prj["namespace"])
}

func GetConfigmapDetails(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetConfigmap(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeleteConfigmap(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteConfigmap(w, r, id["cluster_id"], prj["namespace"], name["name"])

}

//Storage Class

func CreateStorageClass(w http.ResponseWriter, r *http.Request) {
	var sc Models.StorageClass
	dao.CreateStorageClass(w, r, sc)
}

func ListStorageClass(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListSC(w, r, id["cluster_id"])
}

func GetDetailsOfStorageClass(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetStorageClassDetails(w, r, id["cluster_id"], name["name"])
}

func DeleteStorageClass(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteSC(w, r, id["cluster_id"], name["name"])
}

// Persistent Volume

func CreatePersistentVolume(w http.ResponseWriter, r *http.Request) {
	var pv Models.PersistentVolume
	dao.CreatePV(w, r, pv)
}

func ListPersistentVolume(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListPV(w, r, id["cluster_id"])
}

func GetDetailsOfAPersistentVolume(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetPV(w, r, id["cluster_id"], name["name"])
}

func DeletePersistentVolume(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeletePV(w, r, id["cluster_id"], name["name"])
}

// Persistent Volume Claim

func CreatePersistentVolumeClaim(w http.ResponseWriter, r *http.Request) {
	var pvc Models.PersistentVolumeClaim
	namespace := mux.Vars(r)
	dao.CreatePVC(w, r, namespace["namespace"], pvc)
}

func ListPersistentVolumeClaim(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListPVC(w, r, id["cluster_id"], prj["namespace"])
}

func GetPersistentVolumeClaimByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetPVC(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeletePersistentVolumeClaimByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeletePVC(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

// Services

func CreateService(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	var svc Models.Services
	dao.CreateServices(w, r, prj["namespace"], svc)
}

func GetServiceByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetServices(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func ListServiceInNamespace(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListServices(w, r, id["cluster_id"], prj["namespace"])
}

func DeleteServiceByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteServices(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

// //Pods
func CreatePod(w http.ResponseWriter, r *http.Request) {
	var p Models.Pod
	namespace := mux.Vars(r)
	dao.CreatePod(w, r, namespace["namespace"], p)
}

func ListPodsInNamespace(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListPods(w, r, id["cluster_id"], prj["namespace"])
}

func GetPodDetailsByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetPod(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeletePodByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeletePod(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

// Secret

func CreateSecret(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	var secret Models.Secret
	dao.CreateSecret(w, r, prj["namespace"], secret)
}

func ListSecretsDetail(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListSecrets(w, r, id["cluster_id"], prj["namespace"])
}

func GetSecretByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetSecret(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeleteSecretByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteSecret(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

// // Deployment

func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var d Models.Deployment
	prj := mux.Vars(r)
	dao.CreateDeployment(prj["project"], d, w, r)
}

func ListDeploymentsDetail(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListDeployment(w, r, id["cluster_id"], prj["namespace"])
}

func GetDeploymentDetailByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetDeployment(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteDeployment(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

// // Statefulsets

func CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	var sts Models.StatefulSet
	dao.CreateStatefulSet(prj["namespace"], sts, w, r)
}

func ListStatefulSets(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	id := mux.Vars(r)
	dao.ListStatefulSets(w, r, id["cluster_id"], prj["namespace"])
}

func GetStatefulSetByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.GetStatefulSets(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

func DeleteStatefulSetByName(w http.ResponseWriter, r *http.Request) {
	prj := mux.Vars(r)
	name := mux.Vars(r)
	id := mux.Vars(r)
	dao.DeleteStatefulSets(w, r, id["cluster_id"], prj["namespace"], name["name"])
}

//ClusterRole

func CreateClusterRole(w http.ResponseWriter, r *http.Request) {
	var cr Models.ClusterRole
	dao.CreateClusterRole(cr, w, r)
}

func GetClusterRoleDetailsByName(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	name := mux.Vars(r)
	dao.GetClusterRole(id["cluster_id"], name["name"], w, r)
}

func ListClusterRoles(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListClusterRoles(id["cluster_id"], w, r)
}

func DeleteClusterRolesByName(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	name := mux.Vars(r)
	dao.DeleteClusterRole(id["cluster_id"], name["name"], w, r)
}

// Cluster Role Bindings

func CreateClusterRoleBinding(w http.ResponseWriter, r *http.Request) {
	var crb Models.ClusterRoleBinding
	dao.CreateClusterRoleBinding(crb, w, r)
}

func GetClusterRoleBindingDetailsByName(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	name := mux.Vars(r)
	dao.GetClusterRoleBinding(id["cluster_id"], name["name"], w, r)
}

func ListClusterRoleBindings(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	dao.ListClusterRoleBindings(id["cluster_id"], w, r)
}

func DeleteClusterRoleBindingByName(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	name := mux.Vars(r)
	dao.DeleteClusterRoleBinding(id["cluster_id"], name["name"], w, r)
}
