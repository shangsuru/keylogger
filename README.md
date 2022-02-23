# Keylogger

## Server

The server listens at port 2345 for incoming keystrokes sent by the kernel module and stores them in a database, where it can be queried via the API. Start the server:

```
docker-compose up -d --scale worker=5
```

Connect to the database:

```
psql postgresql://postgres:password@<SERVER-IP>:5432/keylogger
```

The server also exposes a RabbitMQ management console at `<SERVER-IP>:8080`.  
Retrieve recorded keystrokes for a given day:

```
curl <SERVER-IP>:5000/api/recordings/2022-02-23
```

## Kernel Module

Follow these steps to install the kernel module onto the linux machine.
It will capture keystrokes and send them to a hard coded (!) SERVER-IP at port 2345.

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
