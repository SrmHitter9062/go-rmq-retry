# go-rmq-retry
This tells the possible solutions for rabbitmq retrying the message for certain time



There are no such feature like retry attempts in RabbitMQ (as well as in AMQP protocol).

Possible solution to implement retry attempts limit behavior:

1. Redeliver message if it was not previously redelivered (check redelivered parameter on basic.deliver method - your library should have some interface for this) and drop it and then catch in dead letter exchange, then process somehow.

2. Each time message cannot be processed publish it again but set or increment/decrement header field, say x-redelivered-count (you can chose any name you like, though). To get control over redeliveries in this case you have to check the field you set whether it reaches some limit (top or bottom - 0 is my choise, a-la ttl in ip header from tcp/ip).

3. Store message unique key in Redis, memcache or other storage, even in mysql alongside with redeliveries count and then on each redelivery increment/decrement this value until it reach the limit.

  3.1) Set requeue count while publishig message and and decrement every requeue time until its reached 0 and then discard message
  3.2) OR set requeue count while requeuing first time ( "abcdsfdsfd-123" => 3) and decrement every requeue time until its reached 0 and then discard message

4.We can have publish_time field in message itself and defined time period (say in T = 5 minutes it should be processed).Before retrying the message by consumer, we can check current time and message published time difference. IF its more than defined time period T then discard message from queue with +ve ack.
