# GRC — Gmail Rule Converter

Gera feeds Atom XML ou XML de importação do Gmail a partir de configurações
YAML, e converte XML do Google Takeout de volta para YAML.

## Funcionalidades

- **Atom XML** — Gera feeds Atom com categorias, autores e propriedades
- **Gmail XML** — Gera XML compatível com importação do Gmail (`--format gmail`)
- **Importar** — Converte XML do Google Takeout para YAML (`--import`)
- **Parser de expressão** — Interpreta sintaxe de busca do Gmail (`from:`, `subject:`, …)
- **Construtor de busca** — Gera links `mail.google.com/#search/…` a partir de critérios
- **Round-trip** — YAML → Gmail XML → YAML sem perda de dados

## Instalar

```bash
make build
sudo cp bin/grc /usr/local/bin/
```

Ou:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Uso

```
grc [opções] <entrada> [saída]
grc --import <gmail.xml> [saída]
```

| Flag              | Descrição                                        |
|-------------------|--------------------------------------------------|
| `--format atom`   | Saída Atom XML (padrão)                          |
| `--format gmail`  | Saída XML de importação do Gmail                 |
| `--import`        | Converte XML do Gmail → YAML                     |
| `--force`         | Sobrescreve arquivo de saída existente           |
| `--verbose`       | Ativa log detalhado                              |
| `--version`       | Exibe versão                                     |
| `--help`          | Exibe ajuda                                      |

### Exemplos

```bash
# Feed Atom (padrão)
grc cfg/example.yaml

# XML de importação do Gmail
grc --format gmail cfg/example-gmail.yaml

# Caminho de saída personalizado
grc --format gmail cfg/example-gmail.yaml meus-filtros.xml

# Importar XML do Google Takeout
grc --import takeout.xml convertido.yaml
```

## YAML — Formato Atom

```yaml
defaults:
  author:
    name: Autor Padrão
    email: autor@exemplo.com
  category:
    term: padrao
    label: Categoria Padrão

filters:
  - name: Newsletter
    expression: subject:newsletter
    category:
      term: newsletter
      label: Newsletters
    properties:
      - name: prioridade
        value: alta
```

Veja `cfg/example.yaml`.

## YAML — Formato Gmail

```yaml
filters:
  - name: Newsletter
    criteria:
      from: newsletter@exemplo.com
    action:
      label: Newsletters
      markRead: true

  - name: Spam
    criteria:
      from: spam@exemplo.com
      hasAttachment: true
    action:
      alwaysSpam: true
```

Campos de critério: `from`, `to`, `subject`, `hasTheWord`, `hasAttachment`,
`is`, `in`, `cc`, `bcc`, `deliveredto`, `list`, `filename`,
`after`, `before`, `older_than`, `newer_than`,
`size`, `larger`, `smaller`, `label`, `category`.

Campos de ação: `archive`, `markRead`, `star`, `delete`, `neverSpam`,
`alwaysSpam`, `forward`, `label`, `smartLabel`, `moveToTrash`.

Veja `cfg/example-gmail.yaml` com 16 filtros de exemplo.

## Compilar

```bash
make build        # bin/grc
make test         # todos os testes
make quality      # fmt + vet + lint + test
make clean        # remove artefatos
```
