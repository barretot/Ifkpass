# 🐞 Debugando Lambda em Go (sem Docker, SAM ou Serverless)

Este guia mostra como **debugar sua função Lambda escrita em Go localmente**, usando apenas ferramentas nativas e leves, como `delve` e `awslambdarpc`.

---

## ✅ 1. Instale os pacotes necessários

### 🐛 Delve (debugger Go)

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 🛠️ awslambdarpc (simula chamada da Lambda)

```bash
go install github.com/blmayer/awslambdarpc@latest
```

> Certifique-se de que `~/go/bin` está no seu `PATH`.

---

## ✅ 2. Configure o VS Code — `.vscode/launch.json`

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "exec",
      "program": "${workspaceFolder}/main",
      "env": {
        "_LAMBDA_SERVER_PORT": "8080"
      },
      "args": []
    }
  ],
  "compounds": [
    {
      "name": "build and debug",
      "configurations": ["Launch"],
      "preLaunchTask": "build-debug"
    }
  ]
}
```

---

## ✅ 3. Configure as tasks — `.vscode/tasks.json`

```json
{
  "version": "2.0.0",
  "inputs": [
    {
      "id": "json",
      "type": "promptString",
      "description": "Caminho do payload JSON",
      "default": ".vscode/events/create-user.json"
    }
  ],
  "tasks": [
    {
      "label": "build-debug",
      "type": "shell",
      "command": "go build -v -gcflags='all=-N -l' -o main ./cmd/lambda"
    },
    {
      "label": "send-event",
      "type": "shell",
      "command": "awslambdarpc -e ${input:json}",
      "problemMatcher": []
    }
  ]
}
```

---

## ✅ 4. Crie um payload de evento — `.vscode/events/create-user.json`

```json
{
  "resource": "/user",
  "path": "/user",
  "httpMethod": "POST",
  "body": "{\"name\":\"John Doe\",\"email\":\"ruan@example.com\"}",
  "isBase64Encoded": false
}
```

---

## ✅ 5. Executando o fluxo completo

### ▶️ Passo 1: Iniciar o debug da Lambda

1. Vá para a aba **"Run and Debug"** (`Ctrl+Shift+D`)
2. Selecione **"build and debug"**
3. Clique em **▶️ Iniciar**

Isso irá:

* Compilar o binário com suporte a debug
* Iniciar a função Lambda escutando localmente na porta `8080`
* Conectar o VS Code com suporte a breakpoints

---

### 📤 Passo 2: Enviar o evento para a Lambda

1. Pressione `Ctrl+Shift+P` (abrir Command Palette)
2. Digite: `Tasks: Run Task`
3. Escolha a task `send-event`
4. Informe ou confirme o caminho do JSON de evento