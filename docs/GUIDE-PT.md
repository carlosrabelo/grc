# GRC — Guia do Gmail Rule Converter

## Visão Geral

GRC converte definições de filtros em YAML para feeds Atom XML.

## Instalação

```bash
git clone https://github.com/carlosrabelo/grc
cd grc
make build
```

## Linha de Comando

```
grc <comando> [flags] <entrada> [saída]
```

### Comandos

| Comando   | Descrição                        |
|-----------|----------------------------------|
| `atom`    | Gera feed Atom XML a partir de YAML |
| `help`    | Exibe ajuda                      |
| `version` | Exibe versão                     |

### Flags comuns

| Flag             | Descrição                      |
|------------------|--------------------------------|
| `--output PATH`  | Caminho de saída explícito     |
| `--force`        | Sobrescreve arquivo existente  |
| `--verbose`      | Ativa log detalhado            |
| `-h, --help`     | Exibe ajuda                    |

## Configuração YAML

```yaml
filters:
  - name: Newsletter
    expression: subject:(newsletter OR weekly OR digest)
    category:
      term: newsletter
      label: Newsletters
```

## Exemplos

```bash
grc atom docs/example.yaml
```

Veja `docs/example.yaml` para um exemplo completo.

## Desenvolvimento

```bash
make build
make test
make quality
```
