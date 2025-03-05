package ProductService

import (
	lg "module/logger"
)

type Product struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory"`
	Brand       string  `json:"brand"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

var Inventory []Product

func LoadProducts() {
	lg.Log.Info("Entered Load Products function")
	products := []Product{
		{"P12345", "Wireless Mouse", "Ergonomic wireless mouse with adjustable DPI.", "Electronics", "Computer Accessories", "Logitech", 29.99, 150},
		{"P12346", "Bluetooth Headphones", "Over-ear Bluetooth headphones with noise cancellation.", "Electronics", "Audio", "Sony", 99.99, 75},
		{"P12347", "Coffee Maker", "12-cup drip coffee maker with programmable settings.", "Home Appliances", "Kitchen", "Hamilton Beach", 49.99, 30},
		{"P12348", "Electric Kettle", "1.7L electric kettle with quick boiling technology.", "Home Appliances", "Kitchen", "Breville", 39.99, 50},
		{"P12349", "Smartphone", "64GB smartphone with 6.5-inch display and dual cameras.", "Electronics", "Mobile Phones", "Samsung", 299.99, 200},
		{"P12350", "Laptop Sleeve", "Neoprene laptop sleeve with extra padding for protection.", "Electronics", "Computer Accessories", "Targus", 19.99, 120},
		{"P12351", "Smart Watch", "Fitness tracking smart watch with heart rate monitor.", "Electronics", "Wearable", "Fitbit", 149.99, 80},
		{"P12352", "Tablet", "10-inch tablet with 64GB storage and Wi-Fi.", "Electronics", "Mobile Devices", "Apple", 499.99, 60},
		{"P12353", "Bluetooth Speaker", "Portable Bluetooth speaker with 12 hours battery life.", "Electronics", "Audio", "JBL", 69.99, 100},
		{"P12354", "Wireless Charger", "Fast wireless charger for smartphones and other devices.", "Electronics", "Mobile Accessories", "Anker", 25.99, 200},
		{"P12355", "Gaming Headset", "Noise-canceling gaming headset with surround sound.", "Electronics", "Audio", "Razer", 129.99, 40},
		{"P12356", "External Hard Drive", "1TB external hard drive for data storage and backup.", "Electronics", "Storage", "Seagate", 59.99, 150},
		{"P12357", "Action Camera", "Waterproof action camera with 4K video recording.", "Electronics", "Cameras", "GoPro", 249.99, 70},
		{"P12358", "Electric Toothbrush", "Rechargeable electric toothbrush with multiple modes.", "Health & Personal Care", "Personal Care", "Oral-B", 79.99, 90},
		{"P12359", "Hair Dryer", "Professional hair dryer with adjustable heat settings.", "Health & Personal Care", "Personal Care", "Dyson", 299.99, 50},
		{"P12360", "Blender", "High-speed blender with multiple blending options.", "Home Appliances", "Kitchen", "Ninja", 119.99, 80},
		{"P12361", "Air Purifier", "HEPA air purifier for home and office use.", "Home Appliances", "Home & Kitchen", "Honeywell", 199.99, 60},
		{"P12362", "Instant Pot", "Electric pressure cooker with multi-functional cooking features.", "Home Appliances", "Kitchen", "Instant Pot", 89.99, 40},
		{"P12363", "Coffee Grinder", "Electric coffee grinder with adjustable grind size.", "Home Appliances", "Kitchen", "Cuisinart", 49.99, 50},
		{"P12364", "Wireless Keyboard", "Compact wireless keyboard with multi-device support.", "Electronics", "Computer Accessories", "Logitech", 29.99, 130},
		{"P12365", "Wireless Earbuds", "True wireless earbuds with noise isolation.", "Electronics", "Audio", "Apple", 199.99, 120},
		{"P12366", "Smart Thermostat", "Smart thermostat with remote control via app.", "Home Appliances", "Smart Home", "Nest", 129.99, 80},
		{"P12367", "Cordless Vacuum", "Lightweight cordless vacuum cleaner with powerful suction.", "Home Appliances", "Cleaning", "Dyson", 349.99, 45},
		{"P12368", "Microwave Oven", "Compact microwave oven with 10 power settings.", "Home Appliances", "Kitchen", "Panasonic", 79.99, 70},
		{"P12369", "Memory Foam Mattress", "10-inch memory foam mattress with cooling technology.", "Furniture", "Bedroom", "Tempur-Pedic", 599.99, 30},
		{"P12370", "LED Desk Lamp", "Adjustable LED desk lamp with touch controls.", "Home Appliances", "Lighting", "Philips", 39.99, 100},
		{"P12371", "Electric Grill", "Indoor electric grill with non-stick surface.", "Home Appliances", "Kitchen", "George Foreman", 59.99, 90},
		{"P12372", "Smart Light Bulb", "Smart light bulb with adjustable color temperature and remote control.", "Home Appliances", "Smart Home", "Philips Hue", 29.99, 150},
		{"P12373", "Car Vacuum Cleaner", "Portable car vacuum cleaner with powerful suction.", "Automotive", "Cleaning", "Black+Decker", 39.99, 120},
		{"P12374", "Portable Power Bank", "10000mAh portable power bank for charging devices on the go.", "Electronics", "Mobile Accessories", "Anker", 25.99, 200},
		{"P12375", "Digital Camera", "12MP digital camera with zoom lens and video recording.", "Electronics", "Cameras", "Canon", 399.99, 65},
		{"P12376", "Smartphone Case", "Protective case for smartphones with shock absorption.", "Electronics", "Mobile Accessories", "OtterBox", 19.99, 200},
		{"P12377", "Mechanical Keyboard", "RGB mechanical keyboard with customizable keys.", "Electronics", "Computer Accessories", "Corsair", 119.99, 90},
		{"P12378", "Gaming Mouse", "High-precision gaming mouse with customizable buttons.", "Electronics", "Computer Accessories", "Razer", 49.99, 110},
		{"P12379", "Noise-Canceling Headphones", "Over-ear noise-canceling headphones with premium sound.", "Electronics", "Audio", "Bose", 299.99, 70},
		{"P12380", "4K Smart TV", "55-inch 4K UHD smart TV with streaming capabilities.", "Electronics", "Home Entertainment", "Samsung", 699.99, 40},
		{"P12381", "Portable SSD", "1TB portable SSD with fast data transfer speed.", "Electronics", "Storage", "Western Digital", 129.99, 130},
		{"P12382", "Smart Doorbell", "Video doorbell with motion detection and night vision.", "Home Appliances", "Smart Home", "Ring", 99.99, 80},
		{"P12383", "Robot Vacuum", "Self-charging robot vacuum with smart navigation.", "Home Appliances", "Cleaning", "iRobot", 249.99, 50},
		{"P12384", "Gaming Chair", "Ergonomic gaming chair with adjustable lumbar support.", "Furniture", "Office", "Secretlab", 399.99, 35},
		{"P12385", "Smart Display", "8-inch smart display with voice assistant integration.", "Electronics", "Smart Home", "Amazon", 129.99, 90},
		{"P12386", "USB-C Hub", "Multi-port USB-C hub with HDMI and USB 3.0 support.", "Electronics", "Computer Accessories", "Anker", 34.99, 140},
		{"P12387", "Streaming Stick", "4K streaming media player with voice control.", "Electronics", "Home Entertainment", "Roku", 49.99, 150},
		{"P12388", "Smart Lock", "Keyless smart lock with Wi-Fi and app control.", "Home Appliances", "Smart Home", "August", 179.99, 60},
		{"P12389", "Wireless Mesh Router", "Tri-band wireless mesh router for whole-home coverage.", "Electronics", "Networking", "Netgear", 299.99, 50},
		{"P12390", "External Monitor", "27-inch 4K UHD monitor with HDR support.", "Electronics", "Computer Accessories", "Dell", 349.99, 45},
		{"P12391", "Smart Ceiling Fan", "Wi-Fi-enabled smart ceiling fan with voice control.", "Home Appliances", "Smart Home", "Hunter", 199.99, 55},
		{"P12392", "Portable Projector", "Compact 1080p portable projector with built-in speakers.", "Electronics", "Home Entertainment", "Epson", 399.99, 40},
		{"P12393", "Standing Desk", "Adjustable height standing desk with electric controls.", "Furniture", "Office", "FlexiSpot", 299.99, 30},
		{"P12394", "Smart Plugs", "Wi-Fi smart plugs with voice assistant integration.", "Home Appliances", "Smart Home", "TP-Link", 24.99, 200},
		{"P12395", "Dash Cam", "Full HD dash cam with loop recording and night vision.", "Automotive", "Cameras", "Garmin", 129.99, 80},
		{"P12396", "Fitness Tracker", "Slim fitness tracker with heart rate and sleep monitoring.", "Electronics", "Wearable", "Garmin", 99.99, 100},
	}

	for _, k := range products {
		Inventory = append(Inventory, k)
	}

	lg.Log.Info("Products loaded successfully in the Inventory List")
}
