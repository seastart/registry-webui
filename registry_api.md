## [media types](https://docs.docker.com/registry/spec/manifest-v2-2/#media-types)

>application/vnd.docker.distribution.manifest.v1+json: schema1 (existing manifest format)  
>application/vnd.docker.distribution.manifest.v2+json: New image manifest format (schemaVersion = 2)  
>application/vnd.docker.distribution.manifest.list.v2+json: Manifest list, aka “fat manifest”  
>application/vnd.docker.container.image.v1+json: Container config JSON  
>application/vnd.docker.image.rootfs.diff.tar.gzip: “Layer”, as a gzipped tar  
>application/vnd.docker.image.rootfs.foreign.diff.tar.gzip: “Layer”, as a gzipped tar that should never be pushed  
>application/vnd.docker.plugin.v1+json: Plugin config JSON

我们最终要分析定位到的是`application/vnd.docker.container.image.v1+json`，这个是镜像的配置文件，里面包含了镜像的基本信息，比如镜像的大小，镜像的创建时间，镜像的作者，镜像的操作系统，镜像的架构等等。
某名称的仓库repository如nginx可能有n个tag如latest、1.0.1，每个tag又可能支持多架构，比如amd64、arm64等等，所以我们需要分析的是每个tag的每个架构的镜像配置文件。

## [仓库列表repositories](https://docs.docker.com/registry/spec/api/#get-catalog)
```json
{
  "repositories": [
    "seastart/nginx",
    "seastart/registry-webui"
  ]
}
```
## [某仓库的标签列表](https://docs.docker.com/registry/spec/api/#get-tags)
```json
{ "name": "seastart/nginx", "tags": [ "dev", "1.0.1" ] }
```

## 某标签支持的多架构镜像[manifest列表](https://docs.docker.com/registry/spec/api/#get-manifest) [“fat manifest”](https://docs.docker.com/registry/spec/manifest-v2-2/#manifest-list) 
传递tag，如果有多架构镜像
```json
{
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "schemaVersion": 2,
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "digest": "sha256:7d9c86320a913927595cb5f6127116c351b7afdb5f08a0abc13beddf0633f66c",
         "size": 2407,
         "platform": {
            "architecture": "amd64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "digest": "sha256:7dff9e1b0a674ab4d60d8dc2779e9f2926cdec54d4e2af7bc6e447967515bbd6",
         "size": 2407,
         "platform": {
            "architecture": "arm64",
            "os": "linux"
         }
      }
   ]
}
```

如果没有多架构，返回的是v1
```json
{
   "schemaVersion": 1,
   "name": "vcs/log-center",
   "tag": "dev",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"id\":\"a431f80ff654de075e8cd702e4cc537f2ecf7462bdcab302378b9081d9df7994\",\"parent\":\"110e3d459615d7f0281bb63df9c3a6421b8635ef2408a060c516a627d619add3\",\"comment\":\"buildkit.dockerfile.v0\",\"created\":\"2023-02-27T09:40:41.531818453Z\",\"container_config\":{\"Cmd\":[\"ENTRYPOINT [\\\"./main\\\"]\"]},\"throwaway\":true}"
      }
   ],
   "signatures": [
      {
         "header": {
            "jwk": {
               "crv": "P-256",
               "kid": "I6JP:C4IR:JSI2:3AI7:K3XR:FT5H:33UC:GO6U:5JAK:MATR:5WCF:PB4P",
               "kty": "EC",
               "x": "7w0sgemxLCTGDW-VW300LNUs4Tch5AuMZm2Wn9BISg4",
               "y": "kFhpePraOfrWoVKaEh6pPpIc2zTM_z_B69X5YYPlwuc"
            },
            "alg": "ES256"
         },
         "signature": "jnuMdfcyIvbKWi542sjZPU4_SIT3wOMMOkOhSZXkj3LsX398G5qafBOYLWmOMdbVw43ADUy8T_MNXq-WpAdJQA",
         "protected": "eyJmb3JtYXRMZW5ndGgiOjEwNDg0LCJmb3JtYXRUYWlsIjoiQ24wIiwidGltZSI6IjIwMjMtMDQtMTFUMDE6MTc6MzlaIn0"
      }
   ]
}
```

有些情况下会遇到[MANIFEST_UNKNOWN错误](https://nova.moe/docker-attestation/)，好像是buildx的问题，需加上`application/vnd.oci.image.index.v1+json`并过滤unknown architecture，本质原因好像是[Attestation storage](https://docs.docker.com/build/attestations/attestation-storage/#attestation-manifest-descriptor)
```json
{
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:405e9582f8bd7b3166afbfac0756b5888c503815087828f10273953a0551049e",
      "size": 1052,
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:ba5f0278c38f91d57c88a5fc838cf0a265cf41c77ed7175d0e201e738fe301ad",
      "size": 1052,
      "platform": {
        "architecture": "arm64",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:00ce56e7891493a003793e9249ca33db5757a1ea0674f580ac2719bc88907af9",
      "size": 566,
      "annotations": {
        "vnd.docker.reference.digest": "sha256:405e9582f8bd7b3166afbfac0756b5888c503815087828f10273953a0551049e",
        "vnd.docker.reference.type": "attestation-manifest"
      },
      "platform": {
        "architecture": "unknown",
        "os": "unknown"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:43f9f365da57f247c16cbbd359ee53b1121ddfb04a2baa52098b67b678e648c5",
      "size": 566,
      "annotations": {
        "vnd.docker.reference.digest": "sha256:ba5f0278c38f91d57c88a5fc838cf0a265cf41c77ed7175d0e201e738fe301ad",
        "vnd.docker.reference.type": "attestation-manifest"
      },
      "platform": {
        "architecture": "unknown",
        "os": "unknown"
      }
    }
  ]
}
```

## 某架构具体image的[manifest详情](https://docs.docker.com/registry/spec/api/#get-manifest) [image manifest](https://docs.docker.com/registry/spec/manifest-v2-2/#image-manifest-field-descriptions) 
如果上一步返回是多架构，遍历传递每一架构的digest；如果上一步没有多架构，直接传递tag；如果上一步有`MANIFEST_UNKNOWN`错误，加上头`application/vnd.oci.image.manifest.v1+json`
遍历`layers`相加`size`，得到image的大小
```json
{
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "schemaVersion": 2,
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "digest": "sha256:e19a44aafe53c43c592241e6b700937fd0cede16a86ef3cbc8c3b4c759e91795",
      "size": 5385
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "digest": "sha256:63b65145d645c1250c391b2d16ebe53b3747c295ca8ba2fcb6b0cf064a4dc21c",
         "size": 3374446
      },
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "digest": "sha256:852dcea3d31bda910441e4a06650362ce0a34eb3b031f89dee09b7d845709ee6",
         "size": 345448
      }
   ]
}
```

## 获取镜像"application/vnd.docker.container.image.v1+json"的[blob详情](https://docs.docker.com/registry/spec/api/#get-blob)

```json
{
    "architecture": "amd64",
    "config":
    {
        "Env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "TZ=Asia/Shanghai", "WORKDIR=/var/www/registy-webui"],
        "Entrypoint": ["./main"],
        "Cmd": ["--config", "./config/default.yml"],
        "WorkingDir": "/var/www/registy-webui",
        "Labels":
        {
            "description": "docker registy webui",
            "maintainer": "dev@seastart.cn",
            "name": "registy-webui"
        },
        "ArgsEscaped": true,
        "OnBuild": null
    },
    "created": "2023-02-23T08:32:23.210335293Z",
    "history": [
    {
        "created": "2023-02-11T04:46:42.449083344Z",
        "created_by": "/bin/sh -c #(nop) ADD file:40887ab7c06977737e63c215c9bd297c0c74de8d12d16ebdf1c3d40ac392f62d in / "
    },
    {
        "created": "2023-02-23T08:32:23.210335293Z",
        "created_by": "WORKDIR /var/www/registy-webui",
        "comment": "buildkit.dockerfile.v0"
    },
    {
        "created": "2023-02-23T08:32:23.210335293Z",
        "created_by": "ENTRYPOINT [\"./main\"]",
        "comment": "buildkit.dockerfile.v0",
        "empty_layer": true
    },
    {
        "created": "2023-02-23T08:32:23.210335293Z",
        "created_by": "CMD [\"--config\" \"./config/default.yml\"]",
        "comment": "buildkit.dockerfile.v0",
        "empty_layer": true
    }],
    "moby.buildkit.buildinfo.v1": "eyJmcm9udGVuZCI6ImRvY2tlcmZpbGUudjAiLCJzb3VyY2VzIjpbeyJ0eXBlIjoiZG9ja2VyLWltYWdlIiwicmVmIjoiZG9ja2VyLmlvL2xpYnJhcnkvYWxwaW5lOjMiLCJwaW4iOiJzaGEyNTY6Njk2NjVkMDJjYjMyMTkyZTUyZTA3NjQ0ZDc2YmM2ZjI1YWJlYjU0MTBlZGMxYzdhODFhMTBiYTNmMGVmYjkwYSJ9XX0=",
    "os": "linux",
    "rootfs":
    {
        "type": "layers",
        "diff_ids": ["sha256:7cd52847ad775a5ddc4b58326cf884beee34544296402c6292ed76474c686d39", "sha256:e6c82f1d0e7cea960ff5adde208c504aafeb4d8e29a2955d7d6ff4a5d7cfe6b3", "sha256:a2a3d069ef51e9ce838168f08e8bec03dbf320c1e648ea1eab2dca12795b8bec", "sha256:4aca5e952bce55f0442c968c7f9b8ae00cd6278b5f410636899a3156ccafe78f", "sha256:b13f4dfcf5f0f41f60c455b849e70d29b72dee90b07712a7669436a2345eac4e", "sha256:0a448abcb806fdce48509539ff775647ae21cb4e164a3c362822b5ed02c0cdec", "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef", "sha256:8f5099550eb54ef90e272794a17824ba429fcc295fbd5a8234a3d51ac591eb57", "sha256:0a43888dc783aab8b5e836ebd2674fb0609be78f9ac21eaf93dee57e6ec699fd", "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"]
    }
}
```