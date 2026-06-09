package routes

import (
	"gk-capital-backend/config"
	"gk-capital-backend/controllers"
	"gk-capital-backend/middleware"
	"gk-capital-backend/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(
	cfg *config.Config,
	authCtrl *controllers.AuthController,
	serviceCtrl *controllers.ServiceController,
	testimonialCtrl *controllers.TestimonialController,
	faqCtrl *controllers.FAQController,
	homepageCtrl *controllers.HomepageController,
	aboutCtrl *controllers.AboutController,
	contactCtrl *controllers.ContactController,
	loanCtrl *controllers.LoanInquiryController,
	seoCtrl *controllers.SEOController,
	dashboardCtrl *controllers.DashboardController,
	mediaCtrl *controllers.MediaController,
) *gin.Engine {
	r := gin.Default()

	// Middlewares
	r.Use(middleware.CorsMiddleware())

	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Diagnostics check
	r.GET("/api/test-db", func(c *gin.Context) {
		var count int64
		// Check table existence and row count
		errCount := config.DB.Table("seo_pages").Count(&count).Error

		// Try migrating
		var migrateErrStr string
		if err := config.DB.AutoMigrate(&models.SEOPage{}); err != nil {
			migrateErrStr = err.Error()
		}

		// Try executing the seed logic with error capture
		type SeedResult struct {
			PageSlug string `json:"page_slug"`
			Status   string `json:"status"`
			Error    string `json:"error,omitempty"`
		}
		var results []SeedResult

		// Try deleting
		deleteErr := config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SEOPage{}).Error
		var deleteErrStr string
		if deleteErr != nil {
			deleteErrStr = deleteErr.Error()
		}

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
			pCopy := p
			err := config.DB.Create(&pCopy).Error
			status := "success"
			errStr := ""
			if err != nil {
				status = "failed"
				errStr = err.Error()
			}
			results = append(results, SeedResult{
				PageSlug: p.PageSlug,
				Status:   status,
				Error:    errStr,
			})
		}

		c.JSON(200, gin.H{
			"migrate_error": migrateErrStr,
			"count_error":   func() string { if errCount != nil { return errCount.Error() }; return "" }(),
			"initial_count": count,
			"delete_error":  deleteErrStr,
			"seed_results":  results,
		})
	})

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
		}

		// Public Routes (Rate limited for writes)
		publicWrites := api.Group("")
		publicWrites.Use(middleware.RateLimitMiddleware(1.0, 5)) // 1 request per second, burst 5
		{
			publicWrites.POST("/contact", contactCtrl.Submit)
			publicWrites.POST("/loan-inquiry", loanCtrl.Submit)
		}

		api.GET("/services", serviceCtrl.GetAll)
		api.GET("/services/:id", serviceCtrl.GetByID)
		api.GET("/testimonials", testimonialCtrl.GetAll)
		api.GET("/faqs", faqCtrl.GetAll)
		api.GET("/homepage", homepageCtrl.GetAll)
		api.GET("/about", aboutCtrl.GetAll)
		api.GET("/seo", seoCtrl.GetBySlug)

		// Public Media Routes
		api.GET("/media/:id", mediaCtrl.Serve)
		api.GET("/media/name/:name", mediaCtrl.ServeByName)

		// Admin Routes (JWT Protected)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Dashboard
			admin.GET("/dashboard", dashboardCtrl.GetStats)

			// Services
			admin.GET("/services", serviceCtrl.AdminGetAll)
			admin.POST("/services", serviceCtrl.Create)
			admin.PUT("/services/:id", serviceCtrl.Update)
			admin.DELETE("/services/:id", serviceCtrl.Delete)

			// Testimonials
			admin.GET("/testimonials", testimonialCtrl.AdminGetAll)
			admin.POST("/testimonials", testimonialCtrl.Create)
			admin.PUT("/testimonials/:id", testimonialCtrl.Update)
			admin.DELETE("/testimonials/:id", testimonialCtrl.Delete)

			// FAQs
			admin.GET("/faqs", faqCtrl.AdminGetAll)
			admin.POST("/faqs", faqCtrl.Create)
			admin.PUT("/faqs/:id", faqCtrl.Update)
			admin.DELETE("/faqs/:id", faqCtrl.Delete)

			// Homepage Sections
			admin.GET("/homepage", homepageCtrl.AdminGetAll)
			admin.POST("/homepage", homepageCtrl.Create)
			admin.PUT("/homepage/:id", homepageCtrl.Update)
			admin.DELETE("/homepage/:id", homepageCtrl.Delete)

			// About Sections
			admin.GET("/about", aboutCtrl.AdminGetAll)
			admin.POST("/about", aboutCtrl.Create)
			admin.PUT("/about/:id", aboutCtrl.Update)
			admin.DELETE("/about/:id", aboutCtrl.Delete)

			// SEO Pages
			admin.GET("/seo", seoCtrl.GetAll)
			admin.POST("/seo", seoCtrl.Create)
			admin.PUT("/seo/:id", seoCtrl.Update)
			admin.DELETE("/seo/:id", seoCtrl.Delete)

			// Contact Inquiries
			admin.GET("/contact-inquiries", contactCtrl.GetAll)
			admin.PATCH("/contact-inquiries/:id", contactCtrl.UpdateStatus)

			// Loan Inquiries
			admin.GET("/loan-inquiries", loanCtrl.GetAll)
			admin.PATCH("/loan-inquiries/:id", loanCtrl.UpdateStatus)

			// Media Assets
			admin.GET("/media", mediaCtrl.AdminGetAll)
			admin.POST("/media/upload", mediaCtrl.Upload)
			admin.DELETE("/media/:id", mediaCtrl.Delete)
		}
	}

	return r
}
