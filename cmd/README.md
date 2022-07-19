# `cmd`

Here you can find the entrypoint source files for all binaries shinpuru generates and uses. 

- [`shinpuru`](shinpuru/) contains the entrypoint for the main shinpuru server.
- [`cmdman`](cmdman/) contains the entrypoint for the command documentation generation tool.
- [`setup`](setup/) contains a CLI tool to create a pre-configured and ready to deploy docker-compose.yml with simple question promts.

You can play around with these two applications and their flags by using the following command.
```
go run cmd/<directory>/main.go -h
```