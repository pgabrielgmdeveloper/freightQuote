

# Como usar

Para iniciar o projeto siga os passos:

1. Na raiz do projeto execute o comando `docker compose up -d` e aguarde os `containers` do `postgres` e `redis` e `o app frete` iniciarem;

2. Ao rodar dar o `docker compose up -d` as migrations serao executadas automaticamente;

3. caso o container do do `app não inicie automaticamente` pode rodar o comando `docker container start {container-name}`
Com isso o sistema já está pronto para o uso, para testar existe algumas formas:

1. `REST API`:
    - Na raiz do projeto tem um arquivo `requests.http` que pode ser usado com a IDE Goland OU com o vscode com o `plugin` do `vscode` de `client rest` https://github.com/Huachao/vscode-restclient ou https://marketplace.visualstudio.com/items?itemName=humao.rest-client

2. `endpoints criados`;
   - temos 2 end points o primeiro simulate

## Arquitetura do projeto
#### o Projeto utilizar da arquitetura hexal ou port and adpaters
- oque nos facilita a substituição de dependencias com facilidade e a testabilidade do codigo




## Estrutura de Pastas

```
freight_quote/
├── cmd/                 # Ponto de entrada da aplicação (main.go)
├── configs/            # Arquivos de configuração (yaml, env, etc.)
├── internal/           # Código interno não exportável (NUCLEO DO SISTEMA)
│   └── domain/         # Núcleo da aplicação (regras de negócio) ENTIDADES, SERVICES E PORTAS     
└── pkg/
    └── infra/          # Implementações concretas de infraestrutura
        ├── database/   # Adaptadores de banco de dados
        ├── http/       # Handlers HTTP (controllers)
        ├── cache/      # Implementação de cache
        └── clients/    # Clients de serviços externos
```

