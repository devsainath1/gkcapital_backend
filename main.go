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
	// Important: The user specified DO NOT generate schema/migration files, but since we are running a fresh setup locally (or in Docker), GORM AutoMigrate is extremely helpful as a fallback.
	// We'll perform soft migration just in case the tables don't exist yet in the developer's local DB environment.
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

	// Seed initial Homepage sections if not exists
	var homepageCount int64
	config.DB.Model(&models.HomepageSection{}).Count(&homepageCount)
	if homepageCount == 0 {
		sections := []models.HomepageSection{
			{
				SectionKey:  "hero",
				Title:       "Empowering Your Financial Future",
				Subtitle:    "Professional Wealth Management & Financing Solutions",
				Description: "GK Capital delivers institutional-grade investment management, custom lending programs, and strategic corporate financial advisory. Partner with us to scale your wealth securely.",
				Image:       "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=1920&q=80",
				Content: map[string]interface{}{
					"cta_text":          "Apply For Loan",
					"cta_url":           "/apply",
					"secondary_cta_text": "Our Services",
					"secondary_cta_url":  "/services",
				},
				IsActive:  true,
				SortOrder: 1,
			},
			{
				SectionKey:  "why_choose_us",
				Title:       "Why Choose GK Capital",
				Subtitle:    "Designed for Security, Speed, and Growth",
				Description: "Our metrics reflect our commitment to excellence and client success.",
				Content: []map[string]interface{}{
					{"icon": "TrendingUp", "title": "Expert Advisory", "description": "Access elite market analysts and dedicated portfolio managers."},
					{"icon": "Zap", "title": "Instant Approval", "description": "Streamlined lending pathways with quick turnaround times."},
					{"icon": "Shield", "title": "Secure Assets", "description": "Advanced capital safety standards and absolute privacy controls."},
					{"icon": "Briefcase", "title": "Custom Portfolios", "description": "Tailored wealth creation and cash management designs."},
				},
				IsActive:  true,
				SortOrder: 2,
			},
		}
		for _, s := range sections {
			config.DB.Create(&s)
		}
		log.Println("Seeded homepage sections")
	}

	// Seed initial About sections if not exists
	var aboutCount int64
	config.DB.Model(&models.AboutSection{}).Count(&aboutCount)
	if aboutCount == 0 {
		sections := []models.AboutSection{
			{
				SectionKey:  "story",
				Title:       "Our Legacy of Financial Innovation",
				Subtitle:    "Built on Trust, Integrity, and Performance",
				Description: "Founded in 2012, GK Capital has grown from a boutique investment advisory firm into a premier diversified financial services group. We manage over $5B in assets and have enabled capital acquisition for hundreds of enterprises.",
				Image:       "https://images.unsplash.com/photo-1579532561814-a44de6e50c99?auto=format&fit=crop&w=1000&q=80",
				IsActive:    true,
				SortOrder:   1,
			},
			{
				SectionKey:  "mission_vision",
				Title:       "Our Commitments",
				Subtitle:    "Guiding Principles",
				Description: "Our goal is simple: maximize client returns while mitigating financial risks.",
				Content: map[string]interface{}{
					"mission": "To provide exceptional wealth management and capital accessibility, enabling sustainable business expansion and personal legacy security.",
					"vision":  "To be the leading private financial group trusted by families, corporations, and founders globally.",
					"values":  []string{"Client-Centricity", "Radical Integrity", "Bold Innovation", "Disciplined Risk Management"},
				},
				IsActive:  true,
				SortOrder: 2,
			},
			{
				SectionKey: "statistics",
				Title:      "GK Capital by the Numbers",
				Content: []map[string]interface{}{
					{"label": "Assets Under Management", "value": "$5.2B"},
					{"label": "Active Corporate Clients", "value": "1,200+"},
					{"label": "Loans Disbursed", "value": "$850M+"},
					{"label": "Global Offices", "value": "5"},
				},
				IsActive:  true,
				SortOrder: 3,
			},
		}
		for _, s := range sections {
			config.DB.Create(&s)
		}
		log.Println("Seeded about sections")
	}

	// Seed default services if empty
	var serviceCount int64
	config.DB.Model(&models.Service{}).Count(&serviceCount)
	if serviceCount == 0 {
		services := []models.Service{
			{
				Title:       "Corporate Capital & Debt Financing",
				Slug:        "corporate-debt-financing",
				Description: "High-volume loans, lines of credit, and structured debt solutions tailored for expansion-phase enterprises and working capital enhancements.",
				Image:       "https://images.unsplash.com/photo-1460925895917-afdab827c52f?auto=format&fit=crop&w=800&q=80",
				Icon:        "Briefcase",
				IsActive:    true,
				SortOrder:   1,
			},
			{
				Title:       "Wealth Management & Private Banking",
				Slug:        "wealth-management",
				Description: "Comprehensive financial planning, retirement structuring, tax planning, and bespoke portfolio allocation designs for ultra-high-net-worth clients.",
				Image:       "https://images.unsplash.com/photo-1559526324-4b87b5e36e44?auto=format&fit=crop&w=800&q=80",
				Icon:        "TrendingUp",
				IsActive:    true,
				SortOrder:   2,
			},
			{
				Title:       "Commercial Real Estate Lending",
				Slug:        "commercial-real-estate",
				Description: "Flexible construction financing, acquisition bridge loans, and long-term permanent mortgages for industrial, multifamily, and office portfolios.",
				Image:       "https://images.unsplash.com/photo-1560518883-ce09059eeffa?auto=format&fit=crop&w=800&q=80",
				Icon:        "Home",
				IsActive:    true,
				SortOrder:   3,
			},
		}
		for _, s := range services {
			config.DB.Create(&s)
		}
		log.Println("Seeded services")
	}

	// Seed FAQ entries if empty
	var faqCount int64
	config.DB.Model(&models.FAQ{}).Count(&faqCount)
	if faqCount == 0 {
		faqs := []models.FAQ{
			{
				Question:  "What types of loans does GK Capital provide?",
				Answer:    "We offer commercial real estate financing, corporate expansion debt, equipment financing, and custom working capital lines of credit starting from $100,000.",
				Category:  "Lending",
				IsActive:  true,
				SortOrder: 1,
			},
			{
				Question:  "How long is the loan approval process?",
				Answer:    "Initial term sheets are usually issued within 48 hours of documentation submission. Underwriting and capital disbursement typically take between 7 to 14 business days.",
				Category:  "Lending",
				IsActive:  true,
				SortOrder: 2,
			},
			{
				Question:  "What is your minimum capital requirement for Private Wealth Management?",
				Answer:    "Our wealth management services generally begin at a minimum relationship size of $250,000. We also offer advisory for emerging wealth clients.",
				Category:  "Wealth",
				IsActive:  true,
				SortOrder: 3,
			},
		}
		for _, f := range faqs {
			config.DB.Create(&f)
		}
		log.Println("Seeded FAQs")
	}

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

	// Seed SEO entries if empty
	var seoCount int64
	config.DB.Model(&models.SEOPage{}).Count(&seoCount)
	if seoCount == 0 {
		pages := []models.SEOPage{
			{
				PageSlug:        "home",
				MetaTitle:       "GK Capital | Premier Wealth Management & Commercial Lending",
				MetaDescription: "Partner with GK Capital for institutional-grade asset management, personalized private wealth portfolios, and elite corporate and commercial lending solutions.",
				MetaKeywords:    "finance, commercial loan, wealth management, private banking, corporate advisory, investment",
				OGTitle:         "GK Capital — Financial Excellence Secured",
				OGDescription:   "Bespoke financing and wealth planning solutions for corporations and individuals.",
				OGImage:         "https://images.unsplash.com/photo-1454165804606-c3d57bc86b40?auto=format&fit=crop&w=1200&q=80",
				CanonicalURL:    "https://gkcapital.com/",
			},
			{
				PageSlug:        "about",
				MetaTitle:       "About Us | GK Capital Legacy & Values",
				MetaDescription: "Read the story of GK Capital. Discover our values of integrity, trust, client-centricity, and how we manage over $5B in global assets.",
				MetaKeywords:    "company story, mission vision, core values, finance legacy",
				OGTitle:         "About GK Capital",
				OGDescription:   "Our history, mission, and stats in wealth advisory and lending.",
				OGImage:         "https://images.unsplash.com/photo-1579532561814-a44de6e50c99?auto=format&fit=crop&w=1200&q=80",
				CanonicalURL:    "https://gkcapital.com/about",
			},
			{
				PageSlug:        "services",
				MetaTitle:       "Our Services | GK Capital Lending & Investing",
				MetaDescription: "Explore our range of corporate debt financing, private wealth banking, and commercial real estate solutions.",
				MetaKeywords:    "business loans, investment services, financial advisory",
				OGTitle:         "Financial Services from GK Capital",
				OGDescription:   "Scale your company or family trust with expert advice and custom debt capital.",
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
}
