# GRC - Gmail Rule Converter

Gera feeds Atom XML a partir de configurações YAML de filtros.

## Instalação

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
sudo cp bin/grc /usr/local/bin/
```

Ou via `go install`:

```bash
go install github.com/carlosrabelo/grc/grc/cmd/grc@latest
```

## Uso

```bash
grc [opções] <entrada.yaml> [saída]
```

### Argumentos

| Argumento    | Descrição                                                 |
|--------------|-----------------------------------------------------------|
| `entrada.yaml` | Caminho do arquivo YAML de configuração (obrigatório)   |
| `saída`       | Caminho do arquivo XML de saída (padrão: `<entrada>.xml`) |

### Opções

| Opção        | Descrição                            |
|--------------|--------------------------------------|
| `--force`    | Sobrescreve arquivo de saída existente |
| `--verbose`  | Ativa log detalhado                  |
| `--version`  | Exibe informações da versão          |
| `--help`     | Exibe esta mensagem de ajuda         |

## Configuração YAML

### Estrutura

```yaml
defaults:
  author:
    name: Autor Padrão
    email: autor@exemplo.com
  category:
    term: padrao
    label: Categoria Padrão

filters:
  - name: Nome do Filtro
    expression: subject:palavra-chave
    category:
      term: slug-da-categoria
      label: Rótulo da Categoria
    author:
      name: Nome do Autor
      email: autor@exemplo.com
    properties:
      - name: nome-da-propriedade
        value: valor-da-propriedade
```

### Campos

| Campo        | Tipo   | Obrigatório | Descrição                             |
|--------------|--------|-------------|---------------------------------------|
| `name`       | string | sim         | Nome do filtro (título da entrada)    |
| `expression` | string | sim         | Expressão de busca                    |
| `category`   | objeto | não         | Categoria Atom (term, scheme, label)  |
| `author`     | objeto | não         | Autor Atom (name, email, uri)         |
| `properties` | lista  | não         | Pares chave-valor no resumo           |

### Valores Padrão

O bloco `defaults` define valores aplicados a qualquer filtro que não os defina explicitamente.

## Exemplos

### Mínimo

```yaml
filters:
  - name: Newsletter
    expression: subject:newsletter
```

### Completo

Veja `cfg/example.yaml` para um exemplo completo com newsletters, chamados de suporte, avisos internos, comunicações com clientes, documentos financeiros e filtros de spam.

## Compilando a Partir do Código Fonte

```bash
make build        # Compila binário em bin/
make test         # Executa todos os testes
make quality      # Formata, verifica e testa
make clean        # Remove artefatos de build
```

Use `make build` ou compile diretamente:

```bash
go build -o bin/grc ./grc/cmd/grc
```

Incorporando versão:

```bash
go build -ldflags="-X github.com/carlosrabelo/grc/grc/internal/app.Version=$(git describe --tags --always --dirty) -X github.com/carlosrabelo/grc/grc/internal/app.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/grc ./grc/cmd/grc
```
