# GoPay Guard

API em Go para receber eventos de pagamento, prevenir webhooks duplicados e calcular risco antifraude com regras de negócio simples e auditáveis.

## Arquitetura

O projeto segue uma organização por domínio dentro de `internal/`:

- `cmd/api`: bootstrap HTTP, configuração de rotas e injeção de dependências.
- `internal/webhook`: validação, idempotência por `event_id`, persistência do evento e disparo da análise de risco.
- `internal/risk`: cálculo do score, classificação `LOW`/`MEDIUM`/`HIGH` e persistência do resultado.
- `internal/risk/rules`: regras puras de score.
- `internal/payments`: listagem paginada de eventos.
- `internal/http/middleware`: autenticação por API Key.
- `internal/database`: conexão PostgreSQL.
- `internal/migrations`: schema SQL inicial.

Fluxo principal:

1. O cliente envia `POST /webhooks/payments` com `X-API-Key`.
2. A API valida o payload.
3. O `event_id` é consultado para evitar processamento duplicado.
4. O evento é salvo em `payment_events`.
5. O serviço de risco calcula score e nível.
6. O resultado é salvo em `payment_risks`.

## Tecnologias

- Go 1.26.2
- Chi Router
- PostgreSQL 16
- Redis 7
- pgx
- Docker Compose
- Air para hot reload no container de desenvolvimento

## Regras de risco

| Regra | Score | Motivo |
| --- | ---: | --- |
| Valor acima de `100000` | +20 | `High amount transaction` |
| IP com mais de 5 eventos | +30 | `ip_repeated` |
| E-mail com 3 ou mais recusas | +25 | `email_many_declines` |
| Cartão usado em outro documento | +40 | `card_many_documents` |
| Device ID com mais de 5 eventos | +30 | `device_id_repeated` |

Classificação:

| Score | Nível |
| ---: | --- |
| `0` a `30` | `LOW` |
| `31` a `70` | `MEDIUM` |
| `71+` | `HIGH` |

## Variáveis de ambiente

Crie o arquivo local a partir do exemplo:

```bash
cp .env.example .env
```

| Variável | Descrição | Exemplo |
| --- | --- | --- |
| `APP_PORT` | Porta HTTP da API | `8080` |
| `DATABASE_URL` | DSN PostgreSQL usado pela aplicação | `postgres://gopay:gopay123@postgres:5432/gopay_db?sslmode=disable` |
| `REDIS_HOST` | Host do Redis | `redis` |
| `REDIS_PORT` | Porta do Redis | `6379` |
| `API_KEY` | Chave exigida no header `X-API-Key` | `change-me` |

## Como executar

Com Docker:

```bash
docker compose up --build
```

Aplicar a migração inicial:

```bash
docker compose exec -T postgres psql -U gopay -d gopay_db < internal/migrations/001_create_payment_tables.sql
```

Rodar os testes:

```bash
go test ./...
```

## Docker

O `docker-compose.yml` sobe três serviços:

- `api`: aplicação Go com Air e volume do projeto montado em `/app`.
- `postgres`: banco principal na porta `5432`.
- `redis`: serviço Redis na porta `6379`.

A API fica disponível em:

```text
http://localhost:8080
```

## Endpoints

### Health check

```http
GET /health
```

Resposta:

```text
GoPay Guard API is running with database
```

### Receber webhook de pagamento

```http
POST /webhooks/payments
X-API-Key: change-me
Content-Type: application/json
```

Request:

```json
{
  "event_id": "8c1a3ba9-3d6f-46bb-bb1d-3df7e1c7e4f1",
  "type": "payment.created",
  "payment_id": "3f8ccf60-6770-4284-a8b6-6bf86f8a9f8f",
  "amount": 120000,
  "status": "declined",
  "currency": "BRL",
  "customer": {
    "email": "cliente@example.com",
    "document": "12345678900"
  },
  "card": {
    "bin": "411111",
    "last4": "1111"
  },
  "ip": "203.0.113.10",
  "device_id": "3e596f5c-1dbb-4216-94f6-3492fd7db2a0"
}
```

Resposta `202 Accepted`:

```json
{
  "status": "webhook received"
}
```

Evento duplicado:

```text
webhook already exists
```

### Consultar risco de um pagamento

```http
GET /payments/{paymentID}/risk
X-API-Key: change-me
```

Resposta:

```json
{
  "payment_id": "3f8ccf60-6770-4284-a8b6-6bf86f8a9f8f",
  "score": 85,
  "level": "HIGH",
  "reasons": ["High amount transaction", "email_many_declines", "card_many_documents"]
}
```

### Listar eventos de pagamento

```http
GET /payments/events?page=1&limit=20
X-API-Key: change-me
```

Resposta:

```json
{
  "data": [
    {
      "id": "",
      "event_id": "8c1a3ba9-3d6f-46bb-bb1d-3df7e1c7e4f1",
      "payment_id": "3f8ccf60-6770-4284-a8b6-6bf86f8a9f8f",
      "type": "payment.created",
      "created_at": "2026-07-17T12:00:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 1,
  "total_pages": 1
}
```

## Roadmap

- Adicionar runner de migrations no startup ou pipeline.
- Criar testes de integração com PostgreSQL.
- Normalizar nomes JSON de resposta para `snake_case` consistente.
- Adicionar observabilidade: logs estruturados, métricas e tracing.
- Evoluir regras de risco para configuração versionada.
- Implementar autenticação por chave rotacionável por cliente.
