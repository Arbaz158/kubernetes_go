package Models

import "k8s.io/apimachinery/pkg/util/intstr"

type Credentials struct {
	ClusterId                string `gorm:"primary_key;not null;unique" json:"ClusterId"`
	Server                   string `json:"Server"`
	CertificateAuthorityData []byte `gorm:"type:longblob" json:"CertificateAuthorityData"`
	ClientCertificateData    []byte `gorm:"type:longblob" json:"ClientCertificateData"`
	Cluster                  string `json:"Cluster"`
	ClientKeyData            []byte `gorm:"type:longblob" json:"ClientKeyData"`
	Token                    string `json:"Token"`
}

type ServiceAccount struct {
	ClusterId  string   `json:"ClusterId"`
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Labels     []Labels `json:"labels"`
	SecretName string   `json:"secretname"`
}

type Labels struct {
	FirstLabel  string `json:"firstlabel"`
	SecondLabel string `json:"secondlabel"`
}

type Namespace struct {
	ClusterId string   `json:"ClusterId"`
	Name      string   `json:"name"`
	Labels    []Labels `json:"labels"`
}

type Configmap struct {
	ClusterId string            `json:"ClusterId"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    []Labels          `json:"labels"`
	Data      map[string]string `json:"data"`
}

type Parameters struct {
	FirstParameter  string `json:"firstparameter"`
	SecondParameter string `json:"secondparameter"`
}

type StorageClass struct {
	ClusterId          string       `json:"ClusterId"`
	Name               string       `json:"name"`
	Namespace          string       `json:"namespace"`
	Labels             []Labels     `json:"labels"`
	VolumeBinidingMode string       `json:"volumebinidingmode"`
	Parameters         []Parameters `json:"parameters"`
	Provisioner        string       `json:"provisioner"`
	ReclaimPolicy      string       `json:"reclaimpolicy"`
	VolumeExpansion    bool         `json:"volumeexpansion"`
}

type HostPath struct {
	MountName string `json:"mountname"`
	Path      string `json:"path"`
	Type      string `json:"type"`
}

type NFS struct {
	MountName string `json:"mountname"`
	Server    string `json:"serverip"`
	Path      string `json:"path"`
}

type AwsEBS struct {
	VolumeId string `json:"volumeid"`
	ReadOnly bool   `json:"readonly"`
	FsType   string `json:"fstype"`
}

type AzureDisk struct {
	Driver       string `json:"driver"`
	ReadOnly     bool   `json:"readonly"`
	VolumeHandle string `json:"volumehandle"`
	FsType       string `json:"fstype"`
}

type GkePD struct {
	Driver          string `json:"driver"`
	VolumeHandle    string `json:"volumehandle"`
	FsType          string `json:"fstype"`
	ReadOnly        bool   `json:"readonly"`
	FirstAttribute  string `json:"firstattribute"`
	SecondAttribute string `json:"secondattribute"`
}

type AzureFile struct {
	Driver          string `json:"driver"`
	VolumeHandle    string `json:"volumehandle"`
	FsType          string `json:"fstype"`
	ReadOnly        bool   `json:"readonly"`
	FirstAttribute  string `json:"firstattribute"`
	SecondAttribute string `json:"secondattribute"`
	SecretName      string `json:"secretname"`
	SecretNamespace string `json:"secretnamespace"`
}

type PersistentVolume struct {
	ClusterId    string    `json:"ClusterId"`
	Name         string    `json:"name"`
	Namespace    string    `json:"namespace"`
	Labels       []Labels  `json:"labels"`
	AccessModes  string    `json:"accessmodes"`
	StorageClass string    `json:"storageclass"`
	Capacity     string    `json:"capacity"`
	MountOptions []string  `json:"mountoptions"`
	PVSource     string    `json:"pvsource"`
	NodeKey      string    `json:"nodekey"`
	NodeOperator string    `json:"NodeOperator"`
	NodeValues   []string  `json:"NodeValues"`
	HostPath     HostPath  `json:"hostpath"`
	NFS          NFS       `json:"nfs"`
	AwsEBS       AwsEBS    `json:"awsebs"`
	AzureDisk    AzureDisk `json:"azuredisk"`
	AzureFile    AzureFile `json:"azurefile"`
	GkePD        GkePD     `json:"gkepd"`
}

type PersistentVolumeClaim struct {
	ClusterId    string   `json:"ClusterId"`
	Name         string   `json:"name"`
	Namespace    string   `json:"Namespace"`
	Labels       []Labels `json:"labels"`
	AccessModes  string   `json:"accessmodes"`
	Capacity     string   `json:"capacity"`
	StorageClass *string  `json:"storageclass"`
	VolumeName   string   `json:"volumename"`
}

type Services struct {
	ClusterId string     `json:"ClusterId"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Labels    []Labels   `json:"labels"`
	Selector  string     `json:"selector"`
	Type      string     `json:"type"`
	NodePort  []NodePort `json:"nodeport"`
	ClusterIp ClusterIp  `json:"clusterip"`
}

type NodePort struct {
	Port       int32              `json:"port"`
	Protocol   string             `json:"protocol"`
	PortName   string             `json:"portname"`
	TargetPort intstr.IntOrString `json:"targetport"`
	NodePort   int32              `json:"nodeport"`
}

type ClusterIp struct {
	Port       int32              `json:"port"`
	Protocol   string             `json:"protocol"`
	PortName   string             `json:"portname"`
	TargetPort intstr.IntOrString `json:"targetport"`
}

type OpaqueSecret struct {
	FirstSecret  string `json:"firstsecret"`
	SecondSecret []byte `json:"secondsecret"`
}

type Secret struct {
	ClusterId    string         `json:"ClusterId"`
	Name         string         `json:"name"`
	Namespace    string         `json:"namespace"`
	Labels       []Labels       `json:"labels"`
	Type         string         `json:"type"`
	DockerConfig string         `json:"dockerconfig"`
	OpaqueSecret []OpaqueSecret `json:"opaquesecret"`
}

type ContainerPort struct {
	Port int32  `json:"port"`
	Name string `json:"name"`
}

type VolumeMount struct {
	MountName string `json:"mountname"`
	MountPath string `json:"mountpath"`
	SubPath   string `json:"subpath"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvFromSecret struct {
	Name       string `json:"name"`
	SecretName string `json:"secretname"`
	SecretKey  string `json:"secretkey"`
}

type EnvFromConfigmap struct {
	Name          string `json:"name"`
	ConfigmapName string `json:"configmapname"`
	ConfigmapKey  string `json:"configmapkey"`
}

type Probes struct {
	Command             []string
	InitialDelaySeconds int32
	TimeoutSeconds      int32
	PeriodSeconds       int32
}

type Container struct {
	ContainerName    string          `json:"containername"`
	Image            string          `json:"image"`
	ContainerPort    []ContainerPort `json:"containerport"`
	VolumeMount      []VolumeMount
	Command          []string           `json:"command"`
	Args             []string           `json:"args"`
	EnvType          string             `json:"envtype"`
	EnvVariable      []EnvVariable      `json:"envvariable"`
	EnvFromSecret    []EnvFromSecret    `json:"envfromsecret"`
	EnvFromConfigmap []EnvFromConfigmap `json:"envfromconfigmap"`
	CpuLimits        int64
	MemoryLimits     string
	CpuRequest       int64
	MemoryRequest    string
	Tty              bool `json:"tty"`
	SecurityContext  SecurityContext
	LivenessProbe    Probes
	ReadinessProbe   Probes
}

type InitContainer struct {
	ContainerName string   `json:"containername"`
	Image         string   `json:"image"`
	Command       []string `json:"command"`
	Args          []string `json:"args"`
	VolumeMount   []VolumeMount
}

type PVC struct {
	ClaimName string `json:"claimname"`
	MountName string `json:"mountname"`
}

type SecretSource struct {
	MountName  string `json:"mountname"`
	SecretName string `json:"secretname"`
	Key        string `json:"keytopath"`
	Path       string `json:"subpath"`
	Mode       *int32 `json:"defaultmode"`
}

type ConfigmapSource struct {
	MountName     string `json:"mountname"`
	ConfigmapName string `json:"configmapname"`
	Key           string `json:"key"`
	Path          string `json:"path"`
	Mode          *int32 `json:"mode"`
}

type VolumeClaimTemplate struct {
	Name             string `json:"name"`
	StorageClassName string `json:"storageclassname"`
	AccessMode       string `json:"accessmode"`
	Capacity         string `json:"capacity"`
}

type Storage struct {
	NFS                 NFS                 `json:"nfs"`
	HostPath            HostPath            `json:"hostpath"`
	PVC                 PVC                 `json:"pvc"`
	SecretSource        SecretSource        `json:"secret"`
	ConfigmapSourceOne  ConfigmapSource     `json:"configmapsourceone"`
	ConfigmapSourceTwo  ConfigmapSource     `json:"configmapsourcetwo"`
	VolumeClaimTemplate VolumeClaimTemplate `json:"volumeclaimtemplate"`
}

type SecurityContext struct {
	RunAsUser  int64 `json:"runasuser"`
	RunAsGroup int64 `json:"runasgroup"`
	FsGroup    int64
	Privileged bool
}

type Pod struct {
	ClusterId       string        `json:"ClusterId"`
	Name            string        `json:"name"`
	Namespace       string        `json:"namespace"`
	Labels          []Labels      `json:"labels"`
	Container       Container     `json:"container"`
	InitContainer   InitContainer `json:"initcontainer"`
	VolumeSource    string        `json:"volumesource"`
	Storage         Storage       `json:"storage"`
	NodeName        string        `json:"nodename"`
	SecurityContext SecurityContext
}

type Deployment struct {
	ClusterId       string   `json:"ClusterId"`
	Name            string   `json:"name"`
	Namespace       string   `json:"namespace"`
	Labels          []Labels `json:"labels"`
	Container       Container
	InitContainer   InitContainer
	Replicas        int32   `json:"replicas"`
	VolumeSource    string  `json:"volumesource"`
	Storage         Storage `json:"storage"`
	NodeName        string  `json:"nodename"`
	SecurityContext SecurityContext
}

type StatefulSet struct {
	ClusterId       string    `json:"ClusterId"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	Labels          []Labels  `json:"labels"`
	Replicas        int32     `json:"replicas"`
	ServiceName     string    `json:"servicename"`
	Container       Container `json:"container"`
	VolumeSource    string    `json:"volumesource"`
	Storage         Storage   `json:"storage"`
	NodeName        string    `json:"nodename"`
	InitContainer   InitContainer
	SecurityContext SecurityContext
}

//3. Cluster role

type RoleRules struct {
	ApiGroup  []string `json:"apigroup"`
	Resources []string `json:"resources"`
	Verbs     []string `json:"verbs"`
}

type ClusterRole struct {
	ClusterId string `json:"ClusterId"`
	Name      string `json:"name"`
	Labels    []Labels
	RoleRules []RoleRules `json:"rolerules"`
}

type Subject struct {
	ApiGroup string `json:"apigroup"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
}

type RoleRef struct {
	ApiGroup string `json:"apigroup"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
}

type ClusterRoleBinding struct {
	ClusterId string    `json:"ClusterId"`
	Name      string    `json:"name"`
	Labels    []Labels  `json:"labels"`
	Subject   []Subject `json:"subjects"`
	RoleRef   RoleRef   `json:"roleref"`
}
