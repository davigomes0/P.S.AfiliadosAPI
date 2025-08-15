# Sistema de Notificação de Conversões de Afiliados

## Introdução
 A solução foi projetada para lidar com requisições HTTP, garantir a autenticidade dos parceiros, evitar duplicidade de dados e persistir as informações em um banco de dados MySQL.

---

## Tecnologias Utilizadas
* **Go (Golang):** A linguagem de programação principal, escolhida por sua performance, concorrência nativa e tipagem forte.
* **Gin Gonic:** Um framework web para Go, utilizado para facilitar a criação de rotas, middlewares e o gerenciamento de requisições HTTP.
* **MySQL:** O banco de dados relacional escolhido para a persistência dos dados, devido à sua confiabilidade e robustez.
* **Docker:** Utilizado para criar um ambiente isolado e de fácil configuração para o servidor MySQL.
* **Git & GitHub:** Para controle de versão e hospedagem do código.

---

## Decisões de Design e Solução
### 1. Autenticação
Para garantir que apenas parceiros confiáveis possam enviar notificações, a API utiliza **chaves de API (API Keys)**. Cada parceiro recebe uma chave única que deve ser enviada no cabeçalho HTTP `X-API-Key` em cada requisição `POST`.

* A API consulta a tabela `partners` no banco de dados para verificar se a chave fornecida é válida.
* Em caso de chave ausente ou inválida, a API retorna um status `401 Unauthorized`.

### 2. Idempotência
Para prevenir que notificações duplicadas (devido a falhas de rede, por exemplo) causem registros repetidos, a solução garante a **idempotência** usando o identificador de transação (`transaction_id`).

* Antes de inserir uma nova conversão, a API faz uma consulta ao banco de dados para verificar a existência de um registro com o mesmo `transaction_id`.
* Se a transação já existir, a API retorna um status `200 OK` (como se a requisição tivesse sido processada com sucesso) sem criar um novo registro.
* A coluna `transaction_id` na tabela `conversions` possui um índice **`UNIQUE`**, que serve como uma camada extra de segurança para garantir a unicidade no nível do banco de dados.

### 3. Persistência de Dados
O banco de dados MySQL foi modelado com um esquema simples e eficiente:

* **Tabela `partners`:** Armazena os parceiros de negócios e suas chaves de API.
* **Tabela `conversions`:** Armazena os registros de conversão. Possui uma chave estrangeira (`partner_id`) que a liga ao parceiro correspondente na tabela `partners`. O valor da venda (`amount`) é armazenado como um tipo **`DECIMAL`** para garantir precisão monetária.

---

## Como Rodar o Projeto
Siga estes passos para configurar e executar a aplicação em seu ambiente local.

### Pré-requisitos
* [**Go**](https://golang.org/dl/)
* [**Docker Desktop**](https://www.docker.com/products/docker-desktop/)
* [**Git**](https://git-scm.com/downloads/)

### Passos
1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/davigomes0/P.S.AfiliadosAPI.git](https://github.com/davigomes0/P.S.AfiliadosAPI.git)
    cd P.S.AfiliadosAPI
    ```

2.  **Configure o banco de dados com Docker:**
    ```bash
    docker run --name afiliados-mysql -e MYSQL_ROOT_PASSWORD=davi1234 -e MYSQL_DATABASE=afiliados_db -p 3306:3306 -d mysql:8.0
    ```

3.  **Execute as migrações SQL:**
    ```bash
    # No Git Bash ou terminal Linux/macOS
    docker exec -i afiliados-mysql mysql -u root -pdavi1234 afiliados_db < internal/database/migrations/tables.sql
    ```
    *(Se estiver usando PowerShell, use o comando alternativo: `type "internal/database/migrations/tables.sql" | docker exec -i afiliados-mysql mysql -u root -pdavi1234 afiliados_db`)*

4.  **Inicie o servidor Go:**
    ```bash
    go run cmd/main.go
    ```

---

## Exemplos de Requisição
Com o servidor rodando, abra um novo terminal (preferencialmente **Git Bash**) e use os comandos `curl` para testar os endpoints.

### Notificação de Conversão Válida
Esta requisição registrará uma nova conversão no banco de dados.
```bash
curl -X POST http://localhost:8080/api/v1/conversions \
-H "Content-Type: application/json" \
-H "X-API-Key: davi-chave-secreta-1234" \
-d '{
  "transaction_id": "tx-12345678",
  "amount": 99.99
}'

### Tentativa de Duplicidade (Teste de Idempotência)

Execute a requisição acima novamente. A API deve retornar um status 200 OK sem criar uma nova entrada.

### Falha de Autenticação

Esta requisição simula um acesso não autorizado, e a API deve retornar um erro.


curl -X POST http://localhost:8080/api/v1/conversions \
-H "Content-Type: application/json" \
-d '{
  "transaction_id": "tx-98765432",
  "amount": 50.00
}'