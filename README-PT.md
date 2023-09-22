# GRC — Gmail Rule Converter

Gera feeds Atom XML ou XML de importação do Gmail a partir de configurações
YAML, e converte XML do Google Takeout de volta para YAML.

## Funcionalidades

- **Atom XML** — Gera feeds Atom com categorias, autores e propriedades
- **Gmail XML** — Gera XML compatível com importação do Gmail via `grc gmail`
- **Importar** — Converte XML do Google Takeout para YAML via `grc import`
- **Parser de expressão** — Interpreta sintaxe de busca do Gmail (`from:`, `subject:`, `has:attachment`, …)
- **Construtor de busca** — Gera links `mail.google.com/#search/…` a partir de critérios
- **Round-trip** — YAML → Gmail XML → YAML sem perda de dados

## Instalação

### Compilar a partir do código-fonte

```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

Instalar em `~/.local/bin` (sem precisar de root):

```bash
make install
```

### Usando Go Install

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Uso

```
grc <comando> [flags] <entrada> [saída]

Comandos:
  atom      Gera feed Atom XML a partir de YAML
  gmail     Gera XML de importação do Gmail a partir de YAML
  import    Importa XML de filtros do Gmail de volta para YAML
  help      Exibe ajuda (use "grc help <comando>" para detalhes)
  version   Exibe a versão

Flags comuns:
  --output PATH   Caminho de saída explícito
  --force         Sobrescreve arquivo de saída existente
  --verbose       Ativa log detalhado
  -h, --help      Exibe ajuda do comando atual
```

### Caminhos de saída padrão

| Comando  | Extensão padrão  | Exemplo                                |
|----------|------------------|----------------------------------------|
| `atom`   | `.xml`           | `cfg/example.yaml` → `example.xml`     |
| `gmail`  | `.gmail.xml`     | `cfg/example.yaml` → `example.gmail.xml` |
| `import` | `.yaml`          | `takeout.xml` → `takeout.yaml`         |

### Exemplos

```bash
# Feed Atom
grc atom cfg/example.yaml

# XML de importação do Gmail
grc gmail cfg/example-gmail.yaml

# Caminho de saída personalizado (posicional)
grc gmail cfg/example-gmail.yaml meus-filtros.xml

# O mesmo, via flag --output
grc gmail --output meus-filtros.xml cfg/example-gmail.yaml

# Importar XML do Google Takeout
grc import takeout.xml convertido.yaml

# Sobrescrever saída existente
grc gmail --force cfg/example-gmail.yaml
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

## Desenvolvimento

```bash
make build      # Compila o binário para bin/grc
make test       # Executa todos os testes
make quality    # fmt + vet + lint + test
make install    # Instala em ~/.local/bin
make clean      # Remove artefatos de build
```

## Licença

Este projeto está licenciado sob a Licença MIT — veja [LICENSE](LICENSE) para detalhes.
