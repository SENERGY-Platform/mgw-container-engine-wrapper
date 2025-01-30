mgw-container-engine-wrapper
=======

![Image](https://img.shields.io/github/v/tag/SENERGY-Platform/mgw-container-engine-wrapper?filter=v%2A&label=release)

![Image](https://img.shields.io/github/v/tag/SENERGY-Platform/mgw-container-engine-wrapper?filter=lib%2A&label=latest)

![Image](https://img.shields.io/github/v/tag/SENERGY-Platform/mgw-container-engine-wrapper?filter=client%2A&label=latest)

Generate swagger docs:

    swag init -g routes.go -o handler/http_hdl/swagger_docs -dir handler/http_hdl/standard,handler/http_hdl/shared --parseDependency --instanceName standard
    swag init -g routes.go -o handler/http_hdl/swagger_docs -dir handler/http_hdl/restricted,handler/http_hdl/shared --parseDependency --instanceName restricted