# GRC - Gmail Rules Creator

Uma ferramenta de linha de comando para gerar arquivos XML de configuração de filtros Gmail a partir de arquivos de configuração YAML. Simplifica o gerenciamento e personalização de filtros de email Gmail com validação e suporte a valores padrão.

## Funcionalidades

- Configuração YAML: Defina filtros usando sintaxe YAML limpa e legível
- Critérios Abrangentes: Suporte para todos os critérios de filtro Gmail (from, to, subject, query, anexos, etc.)
- Ações Completas: Gama completa de ações Gmail (arquivar, marcar como lido, estrela, encaminhar, lixeira, labels, smart labels)
- Valores Padrão: Aplica automaticamente valores padrão de ações booleanas quando omitidas
- Validação: Garante dados do autor, ao menos um filtro e que cada filtro tenha critérios e ações
- Geração XML: Produz XML formatado corretamente compatível com importação de filtros Gmail
- Log Detalhado: Log detalhado opcional para depuração e monitoramento

## Estrutura do Projeto
```
grc/
├── core/             # Módulo Go principal e código fonte
│   ├── cmd/
│   │   └── grc/      # Ponto de entrada principal da aplicação CLI
│   ├── internal/
│   │   ├── app/      # Lógica da aplicação e tratamento CLI
│   │   └── rules/    # Lógica principal de filtragem e geração XML
│   ├── go.mod        # Definição do módulo Go
│   ├── go.sum        # Checksums do módulo Go
│   └── Makefile      # Automação de build do core
├── resources/
│   └── example.yaml  # Exemplo de configuração YAML
├── scripts/          # Scripts de instalação e utilitários
├── bin/              # Binários gerados (com .gitkeep)
├── Makefile          # Automação de build raiz (delega para core/)
├── README.md         # Documentação do projeto
└── README-PT.md      # Documentação em português
```

## Configuração de Filtros

### Critérios (Condições)
Cada filtro deve incluir pelo menos um destes critérios:
- `from` - Corresponder ao endereço de email do remetente
- `to` - Corresponder ao endereço de email do destinatário
- `subject` - Corresponder à linha de assunto do email
- `hasTheWord` - Corresponder a emails contendo palavras específicas
- `doesNotHaveTheWord` - Corresponder a emails que NÃO contêm palavras específicas
- `list` - Corresponder a emails de lista de discussão
- `query` - Usar sintaxe de consulta de pesquisa Gmail
- `hasAttachment` - Corresponder a emails com/sem anexos

### Ações
Cada filtro deve incluir pelo menos uma ação:
- `label` - Aplicar um label aos emails correspondentes
- `smartLabel` - Aplicar smart labels Gmail (Importante, Spam, etc.)
- `forwardTo` - Encaminhar emails correspondentes para outro endereço
- Ações booleanas:
  - `shouldArchive` - Pular caixa de entrada (arquivar)
  - `shouldMarkAsRead` - Marcar como lido automaticamente
  - `shouldStar` - Adicionar estrela aos emails correspondentes
  - `shouldNeverSpam` - Nunca enviar para spam
  - `shouldAlwaysMarkAsImportant` - Sempre marcar como importante
  - `shouldNeverMarkAsImportant` - Nunca marcar como importante
  - `shouldTrash` - Deletar emails correspondentes

Ações booleanas herdam padrões da seção `default` quando não especificadas.

## Pré-requisitos
- Go 1.22 ou superior

## Instalação

### Opção 1: Build a partir do Código Fonte
```bash
git clone https://github.com/carlosrabelo/grc.git
cd grc
make build
```

### Opção 2: Instalar Localmente
```bash
make install
```
Isso instala o binário `grc` no seu diretório bin local (`$HOME/.local/bin` para usuários, `/usr/local/bin` para root) após compilar automaticamente.

## Uso

### Uso Básico
```bash
grc [opções] <arquivo_yaml>
```

### Opções
- `-output <arquivo>` - Especificar caminho do arquivo XML de saída (padrão: mesmo que entrada com extensão .xml)
- `-verbose` - Habilitar saída de log detalhada

### Exemplo de Configuração YAML
```yaml
author:
  name: "João Silva"
  email: "joao.silva@empresa.com"

default:
  shouldArchive: true
  shouldMarkAsRead: false
  shouldStar: false
  shouldNeverSpam: true
  shouldAlwaysMarkAsImportant: false
  shouldNeverMarkAsImportant: false
  shouldTrash: false

filters:
  - from: "info@newsletter.shopee.com.br"
    label: "@SaneLater"
  - to: "suporte@exemplo.com"
    subject: "[Ticket]"
    hasAttachment: true
    label: "@Suporte"
    shouldMarkAsRead: true
    shouldStar: true
  - query: "list:anuncios.exemplo.com"
    label: "@Anuncios"
    forwardTo: "arquivo@exemplo.com"
    shouldArchive: false
    shouldAlwaysMarkAsImportant: true
```

### Exemplos
```bash
# Gerar XML a partir de configuração YAML
grc resources/example.yaml

# Especificar arquivo de saída customizado
grc -output meus-filtros.xml resources/example.yaml

# Habilitar log detalhado
grc -verbose resources/example.yaml
```

## Desenvolvimento

### Targets Make Disponíveis
```bash
# Nível raiz (delega para core/)
make help          # Mostrar targets disponíveis (padrão)
make build         # Fazer build do binário
make test          # Executar suíte de testes
make run           # Executar a aplicação
make clean         # Remover artefatos de build
make install       # Instalar binário localmente
make uninstall     # Remover binário instalado

# Acesso direto ao core/
cd core && make help  # Mostrar targets específicos do core
cd core && make lint   # Executar linter (golangci-lint)
```

### Trabalhando com o Módulo Core
O projeto usa uma estrutura de diretório core/ com seu próprio módulo Go:
- Módulo: `github.com/carlosrabelo/grc/core`
- Todo desenvolvimento Go acontece dentro do diretório core/
- O Makefile raiz delega todos comandos relacionados ao Go para core/Makefile
- Executar `make` sem argumentos mostra a ajuda por padrão

## Licença
Este projeto é licenciado sob a Licença MIT. Veja o arquivo LICENSE para mais detalhes.
