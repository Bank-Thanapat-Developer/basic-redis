package main

import (
	"log"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/config"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// ---- Seed RefItemType ----
	refItemTypes := []entities.RefItemType{
		{Name: "อาหาร", CreatedAt: time.Now()},
		{Name: "เครื่องดื่ม", CreatedAt: time.Now()},
		{Name: "ของใช้", CreatedAt: time.Now()},
		{Name: "อิเล็กทรอนิกส์", CreatedAt: time.Now()},
	}

	for i := range refItemTypes {
		result := db.Where("name = ?", refItemTypes[i].Name).FirstOrCreate(&refItemTypes[i])
		if result.Error != nil {
			log.Fatalf("failed to seed ref_item_type %q: %v", refItemTypes[i].Name, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("[ref_item_types] created: id=%d name=%s", refItemTypes[i].ID, refItemTypes[i].Name)
		} else {
			log.Printf("[ref_item_types] already exists: id=%d name=%s", refItemTypes[i].ID, refItemTypes[i].Name)
		}
	}

	food := refItemTypes[0].ID
	drink := refItemTypes[1].ID
	supply := refItemTypes[2].ID
	electronic := refItemTypes[3].ID

	// ---- Seed Items (100 รายการ) ----
	items := []entities.Item{
		// อาหาร (30 รายการ)
		{Name: "ข้าวผัด", Price: 60, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ผัดกะเพรา", Price: 55, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ส้มตำ", Price: 50, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ต้มยำกุ้ง", Price: 120, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "แกงเขียวหวาน", Price: 80, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ผัดไทย", Price: 70, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวมันไก่", Price: 65, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ก๋วยเตี๋ยวเรือ", Price: 55, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวหมูแดง", Price: 60, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "หมูกระทะ", Price: 299, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ยำวุ้นเส้น", Price: 55, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ลาบหมู", Price: 65, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "น้ำตกหมู", Price: 65, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ปลาทอด", Price: 90, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ไข่เจียว", Price: 40, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ต้มข่าไก่", Price: 80, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "แกงมัสมั่น", Price: 85, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวซอย", Price: 75, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "สุกี้น้ำ", Price: 150, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "บะหมี่เกี้ยว", Price: 60, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวขาหมู", Price: 70, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "หอยทอด", Price: 80, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ปอเปี๊ยะทอด", Price: 45, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวเหนียวมะม่วง", Price: 55, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "บัวลอย", Price: 35, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "วุ้นกะทิ", Price: 30, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ทับทิมกรอบ", Price: 40, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ข้าวโพดคั่ว", Price: 25, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "ไก่ย่าง", Price: 95, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},
		{Name: "หมูปิ้ง", Price: 15, IsActive: true, RefItemTypeID: food, CreatedAt: time.Now()},

		// เครื่องดื่ม (25 รายการ)
		{Name: "น้ำเปล่า", Price: 10, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "กาแฟดำ", Price: 35, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "กาแฟลาเต้", Price: 55, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "ชาเย็น", Price: 35, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "ชานม", Price: 45, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำส้มคั้น", Price: 40, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำมะพร้าว", Price: 35, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำอ้อย", Price: 20, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "โกโก้ร้อน", Price: 45, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "มอคค่า", Price: 60, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "อเมริกาโน่", Price: 50, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "คาปูชิโน่", Price: 55, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "เอสเพรสโซ่", Price: 40, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "สมูทตี้มะม่วง", Price: 65, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำแตงโม", Price: 35, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำสับปะรด", Price: 35, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "ไมโลเย็น", Price: 40, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "โอวัลตินร้อน", Price: 40, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำขิง", Price: 30, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำตะไคร้", Price: 30, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำใบเตย", Price: 25, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "เบียร์สิงห์", Price: 65, IsActive: false, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "เบียร์ช้าง", Price: 60, IsActive: false, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "โซดา", Price: 15, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},
		{Name: "น้ำมะนาว", Price: 25, IsActive: true, RefItemTypeID: drink, CreatedAt: time.Now()},

		// ของใช้ (25 รายการ)
		{Name: "สบู่", Price: 30, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "แชมพู", Price: 89, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ครีมนวดผม", Price: 99, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ยาสีฟัน", Price: 45, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "แปรงสีฟัน", Price: 35, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "โลชั่นบำรุงผิว", Price: 149, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ครีมกันแดด", Price: 199, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "น้ำยาล้างจาน", Price: 55, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ผงซักฟอก", Price: 75, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "น้ำยาปรับผ้านุ่ม", Price: 89, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "กระดาษทิชชู่", Price: 40, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ถุงขยะ", Price: 25, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "น้ำยาถูพื้น", Price: 69, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "สเปรย์ฆ่าแมลง", Price: 120, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ไม้กวาด", Price: 79, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ไม้ถูพื้น", Price: 199, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "กระดาษชำระ", Price: 89, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ผ้าเช็ดหน้า", Price: 49, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "สำลีก้าน", Price: 35, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "แอลกอฮอล์เจล", Price: 55, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "หน้ากากอนามัย", Price: 99, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ไม้จิ้มฟัน", Price: 10, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ยาหม่อง", Price: 45, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "พลาสเตอร์", Price: 35, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},
		{Name: "ยาทากันยุง", Price: 65, IsActive: true, RefItemTypeID: supply, CreatedAt: time.Now()},

		// อิเล็กทรอนิกส์ (20 รายการ)
		{Name: "หูฟัง", Price: 590, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "เมาส์", Price: 350, IsActive: false, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "คีย์บอร์ด", Price: 790, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "เว็บแคม", Price: 1290, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "สายชาร์จ USB-C", Price: 199, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "adapter USB Hub", Price: 590, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "พาวเวอร์แบงค์", Price: 890, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "หลอดไฟ LED", Price: 129, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "ปลั๊กพ่วง", Price: 349, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "สายต่อ HDMI", Price: 249, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "แท่นชาร์จไร้สาย", Price: 690, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "ลำโพงบลูทูธ", Price: 1490, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "ไมโครโฟน USB", Price: 2490, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "กล้องวงจรปิด", Price: 1990, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "เราเตอร์ WiFi", Price: 1290, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "สวิตช์เครือข่าย", Price: 990, IsActive: false, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "แฟลชไดรฟ์ 64GB", Price: 299, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "การ์ด SD 128GB", Price: 490, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "ขาตั้งมือถือ", Price: 159, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
		{Name: "ฟิล์มกันรอยจอ", Price: 99, IsActive: true, RefItemTypeID: electronic, CreatedAt: time.Now()},
	}

	for i := range items {
		result := db.Where("name = ?", items[i].Name).FirstOrCreate(&items[i])
		if result.Error != nil {
			log.Fatalf("failed to seed item %q: %v", items[i].Name, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("[items] created: id=%d name=%s price=%d ref_item_type_id=%d", items[i].ID, items[i].Name, items[i].Price, items[i].RefItemTypeID)
		} else {
			log.Printf("[items] already exists: id=%d name=%s", items[i].ID, items[i].Name)
		}
	}

	log.Println("seed completed successfully")
}
