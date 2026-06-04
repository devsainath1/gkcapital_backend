package main

import (
	"log"

	"gk-capital-backend/config"
	"gk-capital-backend/controllers"
	_ "gk-capital-backend/docs"
	"gk-capital-backend/models"
	"gk-capital-backend/repository"
	"gk-capital-backend/routes"
	"gk-capital-backend/services"
	"gk-capital-backend/utils"

	"gorm.io/gorm"
)

// @title GK Capital API Documentation
// @version 1.0
// @description API server for GK Capital finance website.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to Database
	config.ConnectDatabase(cfg)

	// Auto-Migrate/Verify Schema - since we're using existing tables we don't drop,
	// but we need to ensure the DB connection works and we seed admin / initial data.
	seedData()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB)
	serviceRepo := repository.NewServiceRepository(config.DB)
	testimonialRepo := repository.NewTestimonialRepository(config.DB)
	faqRepo := repository.NewFAQRepository(config.DB)
	homepageRepo := repository.NewHomepageRepository(config.DB)
	aboutRepo := repository.NewAboutRepository(config.DB)
	contactRepo := repository.NewContactRepository(config.DB)
	loanRepo := repository.NewLoanInquiryRepository(config.DB)
	seoRepo := repository.NewSEORepository(config.DB)

	// Initialize Services
	authService := services.NewAuthService(userRepo, cfg)
	serviceService := services.NewServiceService(serviceRepo)
	testimonialService := services.NewTestimonialService(testimonialRepo)
	faqService := services.NewFAQService(faqRepo)
	homepageService := services.NewHomepageService(homepageRepo)
	aboutService := services.NewAboutService(aboutRepo)
	contactService := services.NewContactService(contactRepo)
	loanService := services.NewLoanInquiryService(loanRepo)
	seoService := services.NewSEOService(seoRepo)
	dashboardService := services.NewDashboardService(serviceRepo, testimonialRepo, contactRepo, loanRepo)

	// Initialize Controllers
	authCtrl := controllers.NewAuthController(authService)
	serviceCtrl := controllers.NewServiceController(serviceService)
	testimonialCtrl := controllers.NewTestimonialController(testimonialService)
	faqCtrl := controllers.NewFAQController(faqService)
	homepageCtrl := controllers.NewHomepageController(homepageService)
	aboutCtrl := controllers.NewAboutController(aboutService)
	contactCtrl := controllers.NewContactController(contactService)
	loanCtrl := controllers.NewLoanInquiryController(loanService)
	seoCtrl := controllers.NewSEOController(seoService)
	dashboardCtrl := controllers.NewDashboardController(dashboardService)

	// Setup Router
	r := routes.SetupRouter(
		cfg,
		authCtrl,
		serviceCtrl,
		testimonialCtrl,
		faqCtrl,
		homepageCtrl,
		aboutCtrl,
		contactCtrl,
		loanCtrl,
		seoCtrl,
		dashboardCtrl,
	)

	// Start Server
	log.Printf("Server is running on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func seedData() {
	// Auto migrate tables to ensure they exist (or map schemas if missing)
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.WebsiteSetting{},
		&models.Service{},
		&models.Testimonial{},
		&models.FAQ{},
		&models.ContactInquiry{},
		&models.LoanInquiry{},
		&models.SEOPage{},
		&models.HomepageSection{},
		&models.AboutSection{},
	)
	if err != nil {
		log.Printf("Auto migration warning: %v", err)
	}

	// Seed Super Admin if not exists
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := models.User{
			Name:     "GK Capital Administrator",
			Email:    "admin@gkcapital.com",
			Password: hashedPassword,
			Role:     "SUPER_ADMIN",
			IsActive: true,
		}
		config.DB.Create(&admin)
		log.Println("Seeded admin user: admin@gkcapital.com / admin123")
	}

	// Seed or update Homepage sections
	homepageSections := []models.HomepageSection{
		{
			SectionKey:  "hero",
			Title:       "MSME Loan Consultancy Services",
			Subtitle:    "Your Future Our Funds!",
			Description: "GK Capital is a trusted MSME loan consultancy firm committed to empowering businesses and individuals with the right financial solutions. We work with leading banks and NBFCs to provide quick, hassle-free, and secure funding tailored to your needs.",
			Image:       "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=1920&q=80",
			Content: map[string]interface{}{
				"cta_text":           "Apply For Loan",
				"cta_url":            "/apply",
				"secondary_cta_text": "Our Services",
				"secondary_cta_url":  "/services",
			},
			IsActive:  true,
			SortOrder: 1,
		},
		{
			SectionKey:  "why_choose_us",
			Title:       "Why Choose Us?",
			Subtitle:    "Your Partner in Progress",
			Description: "Our metrics reflect our commitment to excellence and client success.",
			Content: []map[string]interface{}{
				{"icon": "TrendingUp", "title": "Higher Success Rate", "description": "Proven track record of getting loan applications approved swiftly."},
				{"icon": "Users", "title": "100% Assistance", "description": "Complete guidance from application submission to final disbursement."},
				{"icon": "BadgePercent", "title": "Best Interest Rates", "description": "Negotiated terms with top lenders to secure the lowest possible rates."},
				{"icon": "Shield", "title": "Quick & Secure", "description": "Fully secure processing of business documents and private financial data."},
			},
			IsActive:  true,
			SortOrder: 2,
		},
	}
	for _, section := range homepageSections {
		var existing models.HomepageSection
		err := config.DB.Where("section_key = ?", section.SectionKey).First(&existing).Error
		if err != nil {
			config.DB.Create(&section)
		} else {
			existing.Title = section.Title
			existing.Subtitle = section.Subtitle
			existing.Description = section.Description
			existing.Image = section.Image
			existing.Content = section.Content
			existing.IsActive = section.IsActive
			existing.SortOrder = section.SortOrder
			config.DB.Save(&existing)
		}
	}
	log.Println("Seeded/Updated homepage sections")

	// Seed or update About sections
	aboutSections := []models.AboutSection{
		{
			SectionKey:  "story",
			Title:       "About GK Capital",
			Subtitle:    "Your Partner in Progress",
			Description: "GK Capital is a trusted MSME loan consultancy firm committed to empowering businesses and individuals with the right financial solutions. We work with leading banks and NBFCs to provide quick, hassle-free, and secure funding tailored to your needs.",
			Image:       "https://images.unsplash.com/photo-1579532561814-a44de6e50c99?auto=format&fit=crop&w=1000&q=80",
			Content:     map[string]interface{}{},
			IsActive:    true,
			SortOrder:   1,
		},
		{
			SectionKey:  "mission_vision",
			Title:       "Our Commitments",
			Subtitle:    "Guiding Principles",
			Description: "Our goal is simple: bridge the gap between businesses and premier financial institutions.",
			Content: map[string]interface{}{
				"mission": "To empower micro, small, and medium enterprises (MSMEs) with quick, hassle-free, and secure funding solutions that accelerate business growth and sustainability.",
				"vision":  "To be India's most trusted financial consultancy partner, enabling sustainable expansion for businesses and individuals.",
				"values":  []string{"Trust", "Integrity", "Speed", "Customer-Centricity"},
			},
			IsActive:  true,
			SortOrder: 2,
		},
		{
			SectionKey: "statistics",
			Title:      "GK Capital by the Numbers",
			Content: []map[string]interface{}{
				{"label": "Success Rate", "value": "98%"},
				{"label": "Lending Partners", "value": "20+"},
				{"label": "Happy Clients", "value": "5,000+"},
				{"label": "Loans Processed", "value": "₹500Cr+"},
			},
			IsActive:  true,
			SortOrder: 3,
		},
	}
	for _, section := range aboutSections {
		var existing models.AboutSection
		err := config.DB.Where("section_key = ?", section.SectionKey).First(&existing).Error
		if err != nil {
			config.DB.Create(&section)
		} else {
			existing.Title = section.Title
			existing.Subtitle = section.Subtitle
			existing.Description = section.Description
			existing.Image = section.Image
			existing.Content = section.Content
			existing.IsActive = section.IsActive
			existing.SortOrder = section.SortOrder
			config.DB.Save(&existing)
		}
	}
	log.Println("Seeded/Updated about sections")

	// Delete any old services that do not belong to the correct 10 services list
	var allowedSlugs = []string{
		"cgtmse-loan",
		"cc-od-facility",
		"machinery-loan",
		"mudra-loan",
		"unsecured-business-loan",
		"loan-against-property",
		"home-loan",
		"personal-loan",
		"project-loan",
		"red-zone-property-purchase",
	}
	config.DB.Where("slug NOT IN ?", allowedSlugs).Delete(&models.Service{})

	// Seed or update the 10 services
	services := []models.Service{
		{
			Title:       "CGTMSE Loan",
			Slug:        "cgtmse-loan",
			Description: "Collateral-free credit facility up to ₹5 Crores for micro and small enterprises, backed by the government's Credit Guarantee Fund Trust for Micro and Small Enterprises (CGTMSE).",
			Image:       "https://images.unsplash.com/photo-1460925895917-afdab827c52f?auto=format&fit=crop&w=800&q=80",
			Icon:        "Shield",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Title:       "CC & OD",
			Slug:        "cc-od-facility",
			Description: "Cash Credit and Overdraft facilities against inventory, assets, or receivables for flexible, on-demand working capital to manage short-term cash flow gaps efficiently.",
			Image:       "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=800&q=80",
			Icon:        "Briefcase",
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Title:       "Machinery Loan",
			Slug:        "machinery-loan",
			Description: "Asset-backed finance for purchasing, upgrading, or expanding industrial equipment, tools, and advanced machinery with attractive interest rates and flexible repayment tenures.",
			Image:       "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?auto=format&fit=crop&w=800&q=80",
			Icon:        "Settings",
			IsActive:    true,
			SortOrder:   3,
		},
		{
			Title:       "Mudra Loan",
			Slug:        "mudra-loan",
			Description: "Government-backed Pradhan Mantri Mudra Yojana (PMMY) loans up to ₹10 Lakhs under Shishu, Kishore, and Tarun categories to fuel micro-enterprise growth and self-employment.",
			Image:       "https://images.unsplash.com/photo-1559526324-4b87b5e36e44?auto=format&fit=crop&w=800&q=80",
			Icon:        "Award",
			IsActive:    true,
			SortOrder:   4,
		},
		{
			Title:       "Unsecured Business Loan",
			Slug:        "unsecured-business-loan",
			Description: "Quick, collateral-free loans for business expansion, cash flow optimization, or immediate capital needs with minimum documentation and fast processing.",
			Image:       "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?auto=format&fit=crop&w=800&q=80",
			Icon:        "Zap",
			IsActive:    true,
			SortOrder:   5,
		},
		{
			Title:       "Loan Against Property (LAP)",
			Slug:        "loan-against-property",
			Description: "Long-term funding secured by commercial or residential property, offering lower interest rates and flexible tenures for high-value business or personal requirements.",
			Image:       "https://images.unsplash.com/photo-1560518883-ce09059eeffa?auto=format&fit=crop&w=800&q=80",
			Icon:        "Home",
			IsActive:    true,
			SortOrder:   6,
		},
		{
			Title:       "Home Loan",
			Slug:        "home-loan",
			Description: "Tailored housing finance solutions with attractive interest rates, flexible tenures, and quick approval to turn your dream home into reality.",
			Image:       "https://images.unsplash.com/photo-1512917774080-9991f1c4c750?auto=format&fit=crop&w=800&q=80",
			Icon:        "Home",
			IsActive:    true,
			SortOrder:   7,
		},
		{
			Title:       "Personal Loan",
			Slug:        "personal-loan",
			Description: "Convenient personal credit lines for individual financial needs with simple documentation, competitive interest rates, and fast disbursement.",
			Image:       "https://images.unsplash.com/photo-1506126613408-eca07ce68773?auto=format&fit=crop&w=800&q=80",
			Icon:        "User",
			IsActive:    true,
			SortOrder:   8,
		},
		{
			Title:       "Project Loan",
			Slug:        "project-loan",
			Description: "Comprehensive funding for new projects, expansions, and greenfield ventures with structured repayment plans, competitive rates, and end-to-end financial advisory support.",
			Image:       "https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?auto=format&fit=crop&w=800&q=80",
			Icon:        "TrendingUp",
			IsActive:    true,
			SortOrder:   9,
		},
		{
			Title:       "Red Zone Property Purchase",
			Slug:        "red-zone-property-purchase",
			Description: "Specialized financing solutions for purchasing properties in red zone or restricted areas with expert guidance on regulatory compliance, documentation, and bank approvals.",
			Image:       "https://images.unsplash.com/photo-1580587771525-78b9dba3b914?auto=format&fit=crop&w=800&q=80",
			Icon:        "MapPin",
			IsActive:    true,
			SortOrder:   10,
		},
	}
	for _, svc := range services {
		var existing models.Service
		err := config.DB.Where("slug = ?", svc.Slug).First(&existing).Error
		if err != nil {
			config.DB.Create(&svc)
		} else {
			existing.Title = svc.Title
			existing.Description = svc.Description
			existing.Image = svc.Image
			existing.Icon = svc.Icon
			existing.IsActive = svc.IsActive
			existing.SortOrder = svc.SortOrder
			config.DB.Save(&existing)
		}
	}
	log.Println("Seeded/Updated services")

	// Seed or update FAQs
	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.FAQ{})
	faqs := []models.FAQ{
		{
			Question:  "What types of loans does GK Capital consult for?",
			Answer:    "We assist businesses and individuals in securing CGTMSE loans, CC & OD facilities, Machinery loans, Mudra loans, Unsecured Business loans, Loan Against Property (LAP), Home loans, Personal loans, Project loans, and Red Zone Property Purchase financing.",
			Category:  "Lending",
			IsActive:  true,
			SortOrder: 1,
		},
		{
			Question:  "How long does the loan approval and disbursement process take?",
			Answer:    "Initial assessment and bank selection are done within 24-48 hours. The complete bank and NBFC approval to final disbursement typically takes between 7 to 14 business days depending on the loan type and documentation.",
			Category:  "Lending",
			IsActive:  true,
			SortOrder: 2,
		},
		{
			Question:  "What documents are required for an MSME loan?",
			Answer:    "Typically, we need basic KYC documents (PAN, Aadhaar), business registration proof (Udyam registration, GST certificate), last 12 months bank statements, and financial statements (ITRs, Balance Sheets) for higher value loans.",
			Category:  "Documentation",
			IsActive:  true,
			SortOrder: 3,
		},
	}
	for _, f := range faqs {
		config.DB.Create(&f)
	}
	log.Println("Seeded FAQs")

	// Seed Testimonials if empty
	var testCount int64
	config.DB.Model(&models.Testimonial{}).Count(&testCount)
	if testCount == 0 {
		testimonials := []models.Testimonial{
			{
				Name:        "Marcus Sterling",
				Designation: "CEO",
				Company:     "Sterling Logistics Group",
				Content:     "GK Capital secured our Series B bridge loan in record time. Their team understands commercial structures better than any traditional banking institution we worked with.",
				Rating:      5,
				Image:       "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=150&h=150&q=80",
				IsActive:    true,
				SortOrder:   1,
			},
			{
				Name:        "Elena Rostova",
				Designation: "Founder",
				Company:     "Vanguard Tech Lab",
				Content:     "Managing corporate treasury is complicated. GK Capital structured an automated cash-management portfolio that generates outstanding yields while retaining immediate liquidity.",
				Rating:      5,
				Image:       "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&w=150&h=150&q=80",
				IsActive:    true,
				SortOrder:   2,
			},
		}
		for _, t := range testimonials {
			config.DB.Create(&t)
		}
		log.Println("Seeded testimonials")
	}

	// Seed SEO entries
	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SEOPage{})
	pages := []models.SEOPage{
		{
			PageSlug:        "home",
			MetaTitle:       "GK Capital | Premier MSME Loan Consultancy Services",
			MetaDescription: "Partner with GK Capital for quick, secure, and hassle-free MSME loans, machinery loans, working capital funding, and business financing solutions.",
			MetaKeywords:    "msme loan, machinery loan, business funding, working capital, cgtmse loan, cc od facility",
			OGTitle:         "GK Capital — Your Future Our Funds!",
			OGDescription:   "Trusted MSME loan consultancy committed to empowering businesses with leading financial solutions.",
			OGImage:         "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/",
		},
		{
			PageSlug:        "about",
			MetaTitle:       "About Us | GK Capital Legacy & Values",
			MetaDescription: "Read the story of GK Capital. Discover how we partner with leading banks and NBFCs to empower MSMEs and individuals.",
			MetaKeywords:    "about company, loan advisor, financial consulting, partners",
			OGTitle:         "About GK Capital",
			OGDescription:   "Our history, mission, and banking partnerships in loan consulting.",
			OGImage:         "https://images.unsplash.com/photo-1579532561814-a44de6e50c99?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/about",
		},
		{
			PageSlug:        "services",
			MetaTitle:       "Our Loan Services | CGTMSE, Mudra, Business & Home Loans",
			MetaDescription: "Explore our full range of MSME financial services: CGTMSE loans, CC & OD, Machinery loans, Mudra loans, LAP, Home loans, Project loans, and more.",
			MetaKeywords:    "business loans, cgtmse, mudra loan, machinery financing, project loan, red zone property, cc od, personal loan",
			OGTitle:         "MSME Loan Services - GK Capital",
			OGDescription:   "Bespoke financing solutions backed by India's major banks & NBFCs.",
			OGImage:         "https://images.unsplash.com/photo-1460925895917-afdab827c52f?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/services",
		},
		{
			PageSlug:        "calculator",
			MetaTitle:       "EMI Calculator | Loan Interest & Tenure Estimator",
			MetaDescription: "Use GK Capital's smart loan EMI calculator to compute monthly payments, total interest costs, and amortizations in real time.",
			MetaKeywords:    "loan calculator, emi calculator, interest tool, finance planner",
			OGTitle:         "GK Capital EMI Calculator",
			OGDescription:   "Estimate your monthly payments for corporate and commercial loans.",
			OGImage:         "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/calculator",
		},
		{
			PageSlug:        "contact",
			MetaTitle:       "Contact Us | Get in Touch with GK Capital",
			MetaDescription: "Reach our advisory and customer relationship desks for consultations, inquiries, or support.",
			MetaKeywords:    "contact info, support email, phone, location",
			OGTitle:         "Contact GK Capital Team",
			OGDescription:   "Connect with our private lenders and asset advisors.",
			OGImage:         "https://images.unsplash.com/photo-1423666639041-f56000c27a9a?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/contact",
		},
		{
			PageSlug:        "apply",
			MetaTitle:       "Apply for Loan | Secure Capital Lead Form",
			MetaDescription: "Submit your basic credentials and capital requirements to secure an instant consultation and financing assessment.",
			MetaKeywords:    "loan apply, lead form, fast financing, commercial cash",
			OGTitle:         "Apply Now at GK Capital",
			OGDescription:   "Request capital for business expansion or investment portfolios.",
			OGImage:         "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=1200&q=80",
			CanonicalURL:    "https://gkcapital.com/apply",
		},
	}
	for _, p := range pages {
		config.DB.Create(&p)
	}
	log.Println("Seeded SEO pages configurations")
}
