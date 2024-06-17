# GRC — Gmail Rule Converter

Gera feeds Atom XML ou XML de importação do Gmail a partir de configurações YAML, e converte XML do Google Takeout de volta para YAML.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://go.dev/)

## Destaques

- Gera feeds Atom com categorias, autores e propriedades
- Gera XML compatível com importação do Gmail via `grc gmail`
- Converte XML do Google Takeout para YAML via `grc import`
- Interpreta sintaxe de busca do Gmail (`from:`, `subject:`, `has:attachment`, …)
- Gera links `mail.google.com/#search/…` a partir de critérios
- Round-trip YAML → Gmail XML → YAML sem perda de dados

## Pré-requisitos

- **Go 1.22+** — necessário para compilar a partir do código-fonte; [download](https://go.dev/dl/)

## Instalação

### Compilar a partir do código-fonte

```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

Instalar em `~/.local/bin` (padrão), ou em `/usr/local/bin` no sistema (sudo apenas na cópia):

```bash
make install
make install-system
make uninstall
make uninstall-system
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
| `atom`   | `.xml`           | `docs/example.yaml` → `example.xml`     |
| `gmail`  | `.gmail.xml`     | `docs/example.yaml` → `example.gmail.xml` |
| `import` | `.yaml`          | `takeout.xml` → `takeout.yaml`         |

### Exemplos

```bash
# Feed Atom
grc atom docs/example.yaml

# XML de importação do Gmail
grc gmail docs/example-gmail.yaml

# Caminho de saída personalizado (posicional)
grc gmail docs/example-gmail.yaml meus-filtros.xml

# O mesmo, via flag --output
grc gmail --output meus-filtros.xml docs/example-gmail.yaml

# Importar XML do Google Takeout
grc import takeout.xml convertido.yaml

# Sobrescrever saída existente
grc gmail --force docs/example-gmail.yaml
```

## Configuração

### Formato Atom

| Campo        | Tipo    | Descrição                              |
|--------------|---------|----------------------------------------|
| `name`       | string  | Nome do filtro (vira o título da entry)|
| `expression` | string  | Expressão de busca do Gmail            |
| `criteria`   | object  | Critérios estruturados do Gmail        |
| `action`     | object  | **Obrigatório.** Deve incluir `label`  |
| `category`   | object  | Categoria Atom (term, scheme, label)   |
| `author`     | object  | Autor Atom (name, email, uri)          |
| `properties` | array   | Pares chave-valor no summary da entry  |

**Importante**: todo filtro deve especificar `action.label` — a pasta de destino (label do Gmail).
Os filtros são ordenados alfabeticamente por esse label na saída gerada.

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

Veja `docs/example.yaml`.

### Formato Gmail

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

Veja `docs/example-gmail.yaml` com 16 filtros de exemplo.

## Estrutura do Projeto

```
grc/cmd/grc/         # Go entry point
grc/internal/        # Private packages
grc/pkg/types/       # Shared public types
bin/                 # Compiled binaries (git-ignored)
.make/               # Build and install scripts
docs/                # Examples and guides
```

## Desenvolvimento

```bash
make build           # Compila o binário para bin/grc
make test            # Executa todos os testes
make quality         # Formata, executa vet e testes
make install         # Instala o binário em ~/.local/bin
make install-system  # Instala o binário em /usr/local/bin
make uninstall       # Remove de ~/.local/bin
make uninstall-system  # Remove de /usr/local/bin
```

## Contribuição

1. Faça um fork do repositório
2. Crie um branch de feature: `git checkout -b feat/description`
3. Faça o commit com Conventional Commits: `git commit -m "feat: add X"`
4. Envie e abra um pull request

## Licença

Este projeto está licenciado sob a GNU General Public License v3.0 — veja [LICENSE](LICENSE) para detalhes.
