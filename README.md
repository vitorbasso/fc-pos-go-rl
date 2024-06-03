# fc-pos-go-rl
Primeiro desafio pós go expert (rate limiter)

## Requerimentos
  * golang versão 1.22.3 ou superior
  * docker e docker-compose

## Subindo o projeto

 * Para subir um servidor de exemplo com o middleware configurado para um redis local, execute o comando `docker-compose up`
 * O servidor de exemplo ficará disponível no endereço `http://localhost:8080` com apenas um `GET` para a raiz

## Executando os testes

 * Para rodar os testes basta executar o comando `go test ./...` na raiz do projeto

## Funcionamento

 * O rate limiter limita o número de requisições aceitas em um determinado período de tempo (configurável) e quando esse limite é ultrapassado responde com http 429
 * O exemplo presente em `exemples/middleware` demonstra um exemplo de middleware para servidor http utilizando redis como engine de armazenamento.
 * O exemplo possui apenas um `GET http://localhost:8080` com uma resposta de hello world.
 * Ele realiza a limitação baseada em uma de duas possíveis chaves: ou de um token presente no header `API_KEY` ou do IP.

## Configuração do rate limiter

 * Para criar um rate limiter ou um middleware dele, passa-se as opções de configuração como argumento. Segue exemplo:
 ```go
	ratel := ratel.Middleware(
		ratel.WithEnvRules(),
		ratel.WithRedisRequestStore(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		}),
	)
 ```
 O `Middleware` retorna uma função que pode ser usada como middleware

 * Possui uma configuração padrão aplicada a todos os IPs. Por default: 10 req/s com 10s de timeout. É possível alterar esse o padrão com a opção `ratel.WithDefaultRule(item rconfig.ConfigItem)`

 * Possível configurar item por item (com chave token ou ip) por meio da opção `ratel.WithRuleItem(key key.RatelKey, item rconfig.ConfigItem)`

 * Também é possível configurar vários de uma vez por variável de ambiente e utilizando `ratel.WithEnvRules()`. Ele pega duas variáveis de ambiente
 ```.env
RATEL_TOKEN_CONFIG=token,5,1,10 token2,2,1,10 token3,10,2,20
RATEL_IP_CONFIG=::1,5,1,10
```
Onde o primeiro valor corresponde a chave (token ou o ip), o segundo corresponde ao limite de requests, o terceiro ao tempo em que esse limite se aplica e o último a quantidade de tempo a se esperar até voltar a aceitar requests, depois de limitado.
No exemplo acima, temos uma configuração onde
```json
API_KEY=token
reqs/seg = 5
tempo de bloqueio até permitir novas requests = 10s
```

 * É possível utilizar qualquer store como meio de armazenamento, com a opção `WithRequestStore(store request.RequestStore)`. Por padrão, utiliza um armazenamento local em memória (map).

 * Também é possível utilizar o redis como armazenamento, com a opção `WithRedisRequestStore(options *redis.Options)`

O item de configuração `rconfig.ConfigItem` possui os seguintes campos:
```go
	Capacity            int
	TimeWindowSecond    int64
	BlockDurationSecond int64
```
