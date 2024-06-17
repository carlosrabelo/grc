# GRC — Guia do Gmail Rule Converter

## Visão Geral

GRC converte definições de filtros em YAML para:

- **Atom XML** (para leitores de RSS, painéis)
- **XML de importação do Gmail** (compatível com Google Workspace / importação de filtros)
- **Configuração YAML** a partir de XML do Google Takeout (`grc import`)

## Instalação

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
make install          # ~/.local/bin
make package          # .deb com binário em /usr/bin
sudo dpkg -i dist/*.deb
```

Via `go install`:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Linha de Comando

```
grc <comando> [flags] <entrada> [saída]
```

O GRC é organizado em subcomandos. Cada comando recebe um arquivo de entrada
posicional e um caminho de saída opcional.

### Comandos

| Comando   | Descrição                                                         |
|-----------|-------------------------------------------------------------------|
| `atom`    | Gera feed Atom XML a partir de YAML                               |
| `gmail`   | Gera XML de importação do Gmail a partir de YAML                  |
| `import`  | Converte XML de filtro do Gmail (Google Takeout) de volta para YAML |
| `help`    | Exibe ajuda (`grc help <comando>` para ajuda específica)          |
| `version` | Exibe versão e informações de build                               |

### Flags comuns (disponíveis em `atom`, `gmail`, `import`)

| Flag             | Descrição                                                  |
|------------------|------------------------------------------------------------|
| `--output PATH`  | Caminho de saída explícito (sobrescreve o argumento posicional) |
| `--force`        | Sobrescreve arquivo de saída existente                     |
| `--verbose`      | Ativa log detalhado                                        |
| `-h, --help`     | Exibe ajuda do comando atual                               |

### Caminhos de saída padrão

| Comando  | Extensão padrão | Exemplo                                    |
|----------|-----------------|--------------------------------------------|
| `atom`   | `.xml`          | `docs/example.yaml` → `example.xml`         |
| `gmail`  | `.gmail.xml`    | `docs/example.yaml` → `example.gmail.xml`   |
| `import` | `.yaml`         | `takeout.xml` → `takeout.yaml`             |

### Códigos de saída

| Código | Significado     |
|--------|-----------------|
| 0      | Sucesso         |
| 1      | Erro de runtime |
| 2      | Erro de uso     |

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
| `action`     | objeto  | **sim**      | Ações do Gmail a aplicar. **Deve incluir `label`** |
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

> **Importante**: todo filtro **deve** incluir `action.label` — a pasta de destino
> (rótulo do Gmail). Os filtros são ordenados alfabeticamente por este rótulo na
> saída gerada (Atom e Gmail XML).

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

| Campo         | Tipo    | Obrigatório | Descrição                          |
|---------------|---------|-------------|------------------------------------|
| `label`       | string  | **sim**     | Pasta de destino (rótulo do Gmail) |
| `archive`     | bool    | não         | Pular caixa de entrada             |
| `markRead`    | bool    | não         | Marcar como lido                   |
| `star`        | bool    | não         | Destacar mensagem                  |
| `delete`      | bool    | não         | Enviar para lixeira                |
| `neverSpam`   | bool    | não         | Nunca marcar como spam             |
| `alwaysSpam`  | bool    | não         | Sempre marcar como spam            |
| `forward`     | string  | não         | Encaminhar para email              |
| `smartLabel`  | string  | não         | Aplicar categoria de rótulo inteligente |
| `moveToTrash` | bool    | não         | Mover para lixeira                 |

## Criando sua primeira configuração YAML

### 1. Comece com um filtro básico

Crie um arquivo chamado `meus-filtros.yaml`:

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@exemplo.com
    action:
      label: Newsletters
      markRead: true
```

Este filtro encontra emails de `newsletter@exemplo.com`, move para a pasta
`Newsletters` e marca como lido.

### 2. Adicione mais filtros

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@exemplo.com
    action:
      label: Newsletters
      markRead: true

  - name: Faturas
    criteria:
      subject: fatura
      hasAttachment: true
    action:
      label: Faturas
      archive: true
```

### 3. Gere o XML

```bash
# Gerar XML de importação do Gmail
grc gmail meus-filtros.yaml

# Ou gerar um feed Atom
grc atom meus-filtros.yaml
```

### 4. Usando expressões (formato Atom)

Em vez de `criteria`, você pode usar a sintaxe de busca do Gmail:

```yaml
filters:
  - name: Prioridade
    expression: from:chefe@exemplo.com OR subject:(urgente ASAP)
    action:
      label: Importante
      star: true
```

Veja `docs/example.yaml` para uma configuração Atom completa.

### 5. Combinando padrões (formato Atom)

Para feeds Atom, você pode definir autores e categorias padrão:

```yaml
defaults:
  author:
    name: Meu Nome
    email: eu@exemplo.com
  category:
    term: email
    label: Email

filters:
  - name: Newsletter
    expression: subject:newsletter
    action:
      label: Newsletters
```

Cada filtro herda os padrões a menos que sejam substituídos.

### Referência

- `docs/example.yaml` — Configuração Atom completa com 6 filtros
- `docs/example-gmail.yaml` — Configuração Gmail com 16 filtros cobrindo todos os tipos de critério e ação

## Exemplos de uso

### Gerar feed Atom

```bash
grc atom docs/example.yaml
# Produz example.xml
```

### Gerar XML de importação do Gmail

```bash
grc gmail docs/example-gmail.yaml
# Produz example.gmail.xml
```

### Importar filtros do Google Takeout

```bash
grc import takeout-filters.xml
# Produz takeout-filters.yaml
```

### Forçar sobrescrita

```bash
grc gmail --force docs/example-gmail.yaml -o saida.gmail.xml
```

### Modo verboso

```bash
grc gmail --verbose docs/example-gmail.yaml
```

### Round-trip completo

```bash
# YAML → Gmail XML
grc gmail docs/example-gmail.yaml filtros.gmail.xml

# Gmail XML → YAML
grc import filtros.gmail.xml
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

- `docs/example.yaml` — Configuração Atom com 6 filtros baseados em expressão
- `docs/example-gmail.yaml` — Configuração Gmail com 16 filtros cobrindo todos os tipos de critério e ação

## Desenvolvimento

```bash
make build        # Compila para bin/grc
make test         # Executa todos os testes (12 pacotes)
make quality      # Formata, verifica, testa
make clean        # Remove bin/ e .make/stamps
```
