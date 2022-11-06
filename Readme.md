# Users Balance Worker

## _Nice Worker, I guess :smirk:_
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/andrew55516/UsersBalanceWorker)

## Features
___
- Refill users balance
- Add user during refilling, if it doesn't exist
- Check users balance
- Reserve money from users balance on Internal Wallet for some order
- Update order status and transfer money from Internal Wallet to company bill or to users balance according to status 
- Transfer money from user to user
- Provide accounting for needed period
- Provide sorted users operations history

## Tech
___
- Built in [go](https://go.dev/) version 1.18
- Uses [docker](https://www.docker.com/) [_Postgres_](https://www.postgresql.org/) container
- Uses [gin-gonic](https://github.com/gin-gonic/gin)
- Uses [pgx](https://github.com/jackc/pgx)

## API
___
- [Api Documentation](https://app.swaggerhub.com/apis-docs/ANDREYAKSENOV/user-balance-worker_open_api_3_0/1.0#/worker)

## Start
___
- Download and sync dependencies
- Run [docker-compose](https://github.com/andrew55516/UsersBalanceWorker/blob/master/migrations/docker-compose.yaml)
- Create three __Databases__: _users_, _services_, _record_
- Run accorded sql files for each __Database__: [users](https://github.com/andrew55516/UsersBalanceWorker/blob/master/migrations/Users.sql), [services](https://github.com/andrew55516/UsersBalanceWorker/blob/master/migrations/Services.sql), [record](https://github.com/andrew55516/UsersBalanceWorker/blob/master/migrations/Record.sql)
- Enjoy the Worker interaction by [API](https://app.swaggerhub.com/apis-docs/ANDREYAKSENOV/user-balance-worker_open_api_3_0/1.0#/worker)

## Problems I met
___
> ### Do we have store usernames?
> ___
> We need know usernames to make clear comments for users operations, but we add users to our base 
> during refill. We can pass username everytime during refilling, but if we sure that the user already exists in our base, we don't have to.
> So, I decided to keep usernames and made this field optional for _/credit_ request 

> ### How we know if the order was successful or not?
> ___
> Maybe we can set timeout for getting message about order success and set status to __failed__ if time passed.
> But I decided just to wait until we get a message about order status by _/orderStatus_

> ### Can we update order status more than once?
> ___
> Logically we don't. If you have said that's _"ok"_ that's _"ok"_.
> if you have said that's _"failed"_ that's _"failed"_ and there is nothing to change - make new order if you need.

> ### Can we implement orders for unknown services?
> ___
> I guess, we can, so, I did. But in that case we will not be able to make clear comments for users about transactions for that services
> (it looks like __"payment for service #10"__ instead of __"payment for service: ServiceName"__) and for accounting 
> (it looks like __"10,unknown service,1000.000000"__ instead of __"10,ServiceName,1000.000000"__) until them is not added to our base

> ### Can we transfer money from or to non-existing in our base user?
> ___
> As in previous case, I guess, we can, as transfer logic for my Worker is just writing-off given value from one balance
> and refilling another balance by that value. That's why, if we try to transfer money from non-existing user, we will catch 400
> with response ___"msg": "not enough money"___, as balance of non-existing user equals 0. And if we try to transfer money to
> non-existing user, Worker will just add user with this ID to our base and refill his balance by given value, but in that
> case username will be empty until we write it manually, and comments for users operations will look like __"transfer to user #465"__
> instead of __"transfer to user: UserName"__
> 
> Of course, we may simply say "sorry, you can't use some features if you have not refilled balance at least one time yet", but, imho, that sucks

> ### Can we implement orders for non-existing in our base users?
> ___
> The same logic as previous, we will catch 400 with response ___"msg": "not enough money"___

> ### Should we sort history of user operations and do pagination?
> Don't know, but I implemented some basic sort features. What's about pagination, Worker returns list of operations in response 
> to _/history_, so, requesting service may display necessary amount of operations on each page.