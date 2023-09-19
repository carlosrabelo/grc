# GRC — Guia do Gmail Rule Converter

## Visão Geral

GRC converte definições de filtros em YAML para:

- **Atom XML** (para leitores de RSS, painéis)
- **XML de importação do Gmail** (compatível com Google Workspace / importação de filtros)
- **Configuração YAML** a partir de XML do Google Takeout (`--import`)

## Instalação

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
sudo cp bin/grc /usr/local/bin/
```

Via `go install`:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Linha de Comando

```
grc [opções] <entrada> [saída]
grc --import <entrada.xml> [saída]
```

### Argumentos

| Argumento | Descrição                                                         |
|-----------|-------------------------------------------------------------------|
| `entrada` | Arquivo YAML de config ou XML do Gmail (com `--import`)           |
| `saída`   | Caminho de saída (padrão varia conforme o modo — veja tabela)     |

### Caminhos de saída padrão

| Modo                       | Extensão padrão | Exemplo                                  |
|----------------------------|-----------------|------------------------------------------|
| `--format atom` (padrão)   | `.xml`          | `cfg/example.yaml` → `example.xml`       |
| `--format gmail`           | `.gmail.xml`    | `cfg/example.yaml` → `example.gmail.xml` |
| `--import`                 | `.yaml`         | `filters.xml` → `filters.yaml`           |

### Opções

| Flag            | Descrição                                    |
|-----------------|----------------------------------------------|
| `--format atom`   | Saída Atom XML (padrão)                    |
| `--format gmail`  | XML de importação do Gmail                 |
| `--import`      | Importa XML de filtros do Gmail para YAML    |
| `--force`       | Sobrescreve arquivo de saída existente       |
| `--verbose`     | Ativa log detalhado                          |
| `--version`     | Exibe versão e data de compilação            |
| `--help`        | Exibe mensagem de ajuda                      |

### Códigos de saída

| Código | Significado |
|--------|-------------|
| 0      | Sucesso     |
| 1      | Erro        |

## Configuração YAML

Dois formatos de configuração são suportados:

### 1. Formato Atom (baseado em expressão)

Usa expressões de busca no estilo Gmail e metadados específicos do Atom.

```yaml
defaults:
  author:
    name: Autor Padrão
    email: autor@exemplo.com
  category:
    term: email
    scheme: http://exemplo.com/categories
    label: Categoria Email

filters:
  - name: Newsletter
    expression: subject:(newsletter OR weekly OR digest)
    category:
      term: newsletter
      label: Newsletters
    author:
      name: Newsletter Bot
      email: newsletter@exemplo.com
    properties:
      - name: prioridade
        value: alta
```

#### Campos do formato Atom

| Campo        | Tipo    | Obrigatório | Descrição                            |
|--------------|---------|-------------|--------------------------------------|
| `name`       | string  | sim         | Nome do filtro (título da entrada)   |
| `expression` | string  | sim*        | Expressão de busca no estilo Gmail   |
| `criteria`   | objeto  | sim*        | Critérios Gmail estruturados         |
| `action`     | objeto  | não         | Ações do Gmail a aplicar             |
| `category`   | objeto  | não         | Categoria Atom (term, scheme, label) |
| `author`     | objeto  | não         | Autor Atom (name, email, uri)        |
| `properties` | lista   | não         | Pares chave-valor no resumo          |

\* Pelo menos um de `expression` ou `criteria` é obrigatório.

### 2. Formato Gmail (critério + ação)

Usa campos estruturados que mapeiam diretamente para propriedades de filtro do Gmail.

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@exemplo.com
    action:
      label: Newsletters
      markRead: true
```

#### Campos de critério

| Campo          | Tipo    | Descrição                              |
|----------------|---------|----------------------------------------|
| `from`         | string  | Remetente                              |
| `to`           | string  | Destinatário                           |
| `subject`      | string  | Assunto                                |
| `hasTheWord`   | string  | Termos de busca livres                 |
| `hasAttachment` | bool   | Possui anexo                           |
| `is`           | string  | `unread`, `starred`, `snoozed`, `important`, `chat` |
| `in`           | string  | `inbox`, `trash`, `spam`, `anywhere`, `sent`, `drafts` |
| `cc`           | string  | CC                                     |
| `bcc`          | string  | BCC                                    |
| `deliveredto`  | string  | Endereço de entrega                    |
| `list`         | string  | Lista de discussão                     |
| `filename`     | string  | Nome do anexo                          |
| `after`        | string  | Data após (ex.: `2024/01/01`)          |
| `before`       | string  | Data antes                             |
| `older_than`   | string  | Idade relativa (ex.: `30d`, `2y`)      |
| `newer_than`   | string  | Idade relativa (ex.: `7d`)             |
| `size`         | string  | Tamanho exato (ex.: `10M`)             |
| `larger`       | string  | Maior que (ex.: `5M`)                  |
| `smaller`      | string  | Menor que (ex.: `1M`)                  |
| `label`        | string  | Possui rótulo                          |
| `category`     | string  | Categoria (`primary`, `social`, `promotions`, `updates`, `forums`) |

#### Campos de ação

| Campo         | Tipo    | Descrição                          |
|---------------|---------|------------------------------------|
| `archive`     | bool    | Pular caixa de entrada             |
| `markRead`    | bool    | Marcar como lido                   |
| `star`        | bool    | Destacar mensagem                  |
| `delete`      | bool    | Enviar para lixeira                |
| `neverSpam`   | bool    | Nunca marcar como spam             |
| `alwaysSpam`  | bool    | Sempre marcar como spam            |
| `forward`     | string  | Encaminhar para email              |
| `label`       | string  | Aplicar rótulo                     |
| `smartLabel`  | string  | Aplicar categoria de rótulo inteligente |
| `moveToTrash` | bool    | Mover para lixeira                 |

## Exemplos de uso

### Gerar feed Atom

```bash
grc cfg/example.yaml
# Produz example.xml
```

### Gerar XML de importação do Gmail

```bash
grc --format gmail cfg/example-gmail.yaml
# Produz example.gmail.xml
```

### Importar filtros do Google Takeout

```bash
grc --import takeout-filters.xml
# Produz takeout-filters.yaml
```

### Forçar sobrescrita

```bash
grc --format gmail --force cfg/example-gmail.yaml saida.gmail.xml
```

### Modo verboso

```bash
grc --verbose --format gmail cfg/example-gmail.yaml
```

### Round-trip completo

```bash
# YAML → Gmail XML
grc --format gmail cfg/example-gmail.yaml filtros.gmail.xml

# Gmail XML → YAML
grc --import filtros.gmail.yaml
```

## Parser de expressão

Ao usar `expression` (formato Atom), o GRC interpreta a sintaxe de busca do Gmail:

```
from:newsletter@example.com subject:(weekly OR digest) has:attachment
```

Operadores suportados:

| Operador            | Mapeia para                  |
|---------------------|------------------------------|
| `from:`             | Criteria.From                |
| `to:`               | Criteria.To                  |
| `subject:`          | Criteria.Subject             |
| `has:`              | `has:attachment` → HasAttachment=true |
| `is:`               | Criteria.Is                  |
| `in:`               | Criteria.In                  |
| `cc:`               | Criteria.Cc                  |
| `bcc:`              | Criteria.Bcc                 |
| `deliveredto:`      | Criteria.DeliveredTo         |
| `list:`             | Criteria.List                |
| `filename:`         | Criteria.Filename            |
| `after:`            | Criteria.After               |
| `before:`           | Criteria.Before              |
| `older_than:`       | Criteria.OlderThan           |
| `newer_than:`       | Criteria.NewerThan           |
| `size:`             | Criteria.Size                |
| `larger:`           | Criteria.Larger              |
| `smaller:`          | Criteria.Smaller             |
| `label:`            | Criteria.Label               |
| `category:`         | Criteria.Category            |
| Tokens não reconhecidos | Criteria.HasTheWord       |

## Construtor de URL de busca

GRC pode gerar links `mail.google.com` a partir de critérios:

- `query.BuildSearchString(criteria)` → `from:usuario@exemplo.com subject:ola`
- `query.BuildSearchURL(criteria)` → `https://mail.google.com/mail/u/0/#search/from%3Ausuario%40exemplo.com%20subject%3Aola`

## Exemplos

- `cfg/example.yaml` — Configuração Atom com 6 filtros baseados em expressão
- `cfg/example-gmail.yaml` — Configuração Gmail com 16 filtros cobrindo todos os tipos de critério e ação

## Desenvolvimento

```bash
make build        # Compila para bin/grc
make test         # Executa todos os testes (12 pacotes)
make quality      # Formata, verifica, lint, testa
make clean        # Remove bin/ e .make/stamps

# Lint (requer golangci-lint)
make lint
```
