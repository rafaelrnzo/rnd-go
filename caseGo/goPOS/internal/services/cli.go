package services

import (
	"fmt"
	"goPOS/internal/models"
	"goPOS/internal/storage"
	"goPOS/internal/utils"
)

func RunCLI() {
	SeedProduct()
	for {
		fmt.Println()
		fmt.Println("=== POS CLI ===")
		fmt.Println("1) List Products")
		fmt.Println("2) Add Product")
		fmt.Println("3) New Order")
		fmt.Println("4) Report")
		fmt.Println("5) Save Data (JSON)")
		fmt.Println("6) Load Data (JSON)")
		fmt.Println("0) Exit")
		choice, _ := utils.ReadInt("> Choose menu: ")
		fmt.Println()

		switch choice {
		case 1:
			handleListProducts()
		case 2:
			handleAddProduct()
		case 3:
			handleNewOrder()
		case 4:
			handleReport()
		case 5:
			handleSave()
		case 6:
			handleLoad()
		case 0:
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("⚠️  Menu tidak dikenal.")
		}
	}
}

func handleListProducts() {
	products := ListProducts()
	if len(products) == 0 {
		fmt.Println("No products found.")
		return
	}
	fmt.Printf("%-4s %-10s %-24s %s\n", "ID", "SKU", "Name", "Price")
	for _, p := range products {
		fmt.Printf("%-4d %-10s %-24s %s\n", p.ID, p.SKU, p.Name, utils.FormatRupiah(p.Price))
	}
}

func handleAddProduct() {
	sku, _ := utils.ReadLine("SKU: ")
	name, _ := utils.ReadLine("NAME: ")
	price, _ := utils.ReadInt("PRICE (Rp): ")
	if err := AddProduct(sku, name, price); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Added product")
}

func handleNewOrder() {
	order := NewOrder()
	fmt.Printf("Order baru dibuat: #%d (DRAFT)\n", order.ID)
	for {
		fmt.Println()
		fmt.Println("a) Add Item")
		fmt.Println("b) Remove Item")
		fmt.Println("c) Update Qty")
		fmt.Println("d) Show Cart")
		fmt.Println("e) Set Discount")
		fmt.Println("f) Checkout")
		fmt.Println("x) Back")
		in, _ := utils.ReadLine("> Pilih: ")
		switch in {
		case "a":
			handleAddItem(order)
		case "b":
			handleRemoveItem(order)
		case "c":
			handleUpdateQty(order)
		case "d":
			showCart(order)
		case "e":
			handleSetDiscount(order)
		case "f":
			if handleCheckout(order) {
				return
			}
		case "x":
			return
		default:
			fmt.Println("⚠️  Opsi tidak dikenal.")
		}
	}
}

func handleAddItem(order *models.Order) {
	by, _ := utils.ReadLine("by: ")
	switch by {
	case "id":
		id, _ := utils.ReadInt("Product id: ")
		p, err := FindProductByID(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		qty, _ := utils.ReadInt("Quantity: ")
		if err := AddItem(order, p, qty); err != nil {
			fmt.Println(err)
			return
		}
	case "sku":
		sku, _ := utils.ReadLine("SKU: ")
		p, err := FindProductBySKU(sku)
		if err != nil {
			fmt.Println(err)
			return
		}
		qty, _ := utils.ReadInt("Quantity: ")
		if err := AddItem(order, p, qty); err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("Input 'id' or 'sku'")
	}
	fmt.Println("Added item")
	showCart(order)
}

func handleRemoveItem(order *models.Order) {
	id, _ := utils.ReadInt("Product id: ")
	if err := RemoveItem(order, id); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Removed item")
	showCart(order)
}

func handleUpdateQty(order *models.Order) {
	id, _ := utils.ReadInt("Product id: ")
	qty, _ := utils.ReadInt("Quantity: ")
	if err := UpdateQty(order, id, qty); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updated qty")
	showCart(order)
}

func handleSetDiscount(order *models.Order) {
	disc, _ := utils.ReadInt("Discount: ")
	if err := SetDiscount(order, disc); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Set discount")
	showCart(order)
}

func handleCheckout(order *models.Order) bool {
	if len(order.Items) == 0 {
		fmt.Println("No items found.")
		return false
	}
	method, _ := utils.ReadLine("input (cash/ewallet): ")
	switch method {
	case "cash":
		paid, _ := utils.ReadInt("Pay (Rp): ")
		if err := Checkout(order, "cash", paid); err != nil {
			fmt.Println("Error:", err)
			return false
		}
	case "ewallet":
		if err := Checkout(order, "cash", 0); err != nil {
			fmt.Println("Error:", err)
			return false
		}
	default:
		fmt.Println("Input 'cash' or 'ewallet'")
		return false
	}
	printReceipt(order)
	return true
}

func showCart(order *models.Order) {
	fmt.Println("Cart:")
	if len(order.Items) == 0 {
		fmt.Println("No items found.")
	} else {
		fmt.Printf("%-4s %-24s %-8s %-4s %s\n", "ID", "Name", "Price", "Qty", "Line")
		for _, it := range order.Items {
			line := it.Price * it.Qty
			fmt.Printf("%-4d %-24s %-8s %-4d %s\n",
				it.ProductID, it.Name, utils.FormatRupiah(it.Price), it.Qty, utils.FormatRupiah(line))
		}
	}
	fmt.Println("-------------------------------")
	fmt.Println("Subtotal:", utils.FormatRupiah(order.Subtotal))
	fmt.Println("Discount:", utils.FormatRupiah(order.Discount))
	fmt.Println("Total   :", utils.FormatRupiah(order.Total))
}

func printReceipt(order *models.Order) {
	fmt.Println("\n===== RECEIPT =====")
	showCart(order)
	fmt.Println("Paid    :", utils.FormatRupiah(order.Paid))
	fmt.Println("Method  :", order.Method)
	fmt.Println("Change  :", utils.FormatRupiah(order.Change))
	fmt.Println("Status  :", order.Status)
	fmt.Println("===================")
}

func handleReport() {
	r := GenerateReport()
	fmt.Println("=== REPORT (session) ===")
	fmt.Println("Transactions:", r.Transaction)
	fmt.Println("Revenue     :", utils.FormatRupiah(r.Revenue))
}

func handleSave() {
	filename, _ := utils.ReadLine("Filename: ")
	if filename == "" {
		fmt.Println("Filename is empty")
		return
	}
	if err := storage.SaveToFile(filename); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File saved.")
}

func handleLoad() {
	filename, _ := utils.ReadLine("Filename: ")
	if filename == "" {
		fmt.Println("Filename is empty")
		return
	}
	if err := storage.LoadFromFile(filename); err != nil {
		fmt.Println(err)
		return
	}
}
