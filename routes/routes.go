package routes

import (
	"gk-capital-backend/config"
	"gk-capital-backend/controllers"
	"gk-capital-backend/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	userCtrl *controllers.UserController,
	websiteSettingCtrl *controllers.WebsiteSettingController,
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

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
		}

		// Public Routes (Rate limited for writes)
		publicWrites := api.Group("")
		publicWrites.Use(middleware.RateLimitMiddleware(0.05, 2)) // 1 request every 20 seconds, burst 2
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
		api.GET("/settings", websiteSettingCtrl.GetPublic)

		// Public Media Routes
		api.GET("/media/:id", mediaCtrl.Serve)
		api.GET("/media/name/:name", mediaCtrl.ServeByName)

		// Admin Routes (JWT Protected)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Manager Group: MANAGER and SUPER_ADMIN
			managerGroup := admin.Group("")
			managerGroup.Use(middleware.RequireRole("MANAGER", "SUPER_ADMIN"))
			{
				// Dashboard
				managerGroup.GET("/dashboard", dashboardCtrl.GetStats)

				// Contact Inquiries
				managerGroup.GET("/contact-inquiries", contactCtrl.GetAll)
				managerGroup.PATCH("/contact-inquiries/:id", contactCtrl.UpdateStatus)

				// Loan Inquiries
				managerGroup.GET("/loan-inquiries", loanCtrl.GetAll)
				managerGroup.PATCH("/loan-inquiries/:id", loanCtrl.UpdateStatus)
			}

			// Super Admin Group: SUPER_ADMIN only
			superAdminGroup := admin.Group("")
			superAdminGroup.Use(middleware.RequireRole("SUPER_ADMIN"))
			{
				// Services
				superAdminGroup.GET("/services", serviceCtrl.AdminGetAll)
				superAdminGroup.POST("/services", serviceCtrl.Create)
				superAdminGroup.PUT("/services/:id", serviceCtrl.Update)
				superAdminGroup.DELETE("/services/:id", serviceCtrl.Delete)

				// Testimonials
				superAdminGroup.GET("/testimonials", testimonialCtrl.AdminGetAll)
				superAdminGroup.POST("/testimonials", testimonialCtrl.Create)
				superAdminGroup.PUT("/testimonials/:id", testimonialCtrl.Update)
				superAdminGroup.DELETE("/testimonials/:id", testimonialCtrl.Delete)

				// FAQs
				superAdminGroup.GET("/faqs", faqCtrl.AdminGetAll)
				superAdminGroup.POST("/faqs", faqCtrl.Create)
				superAdminGroup.PUT("/faqs/:id", faqCtrl.Update)
				superAdminGroup.DELETE("/faqs/:id", faqCtrl.Delete)

				// Homepage Sections
				superAdminGroup.GET("/homepage", homepageCtrl.AdminGetAll)
				superAdminGroup.POST("/homepage", homepageCtrl.Create)
				superAdminGroup.PUT("/homepage/:id", homepageCtrl.Update)
				superAdminGroup.DELETE("/homepage/:id", homepageCtrl.Delete)

				// About Sections
				superAdminGroup.GET("/about", aboutCtrl.AdminGetAll)
				superAdminGroup.POST("/about", aboutCtrl.Create)
				superAdminGroup.PUT("/about/:id", aboutCtrl.Update)
				superAdminGroup.DELETE("/about/:id", aboutCtrl.Delete)

				// SEO Pages
				superAdminGroup.GET("/seo", seoCtrl.GetAll)
				superAdminGroup.POST("/seo", seoCtrl.Create)
				superAdminGroup.PUT("/seo/:id", seoCtrl.Update)
				superAdminGroup.DELETE("/seo/:id", seoCtrl.Delete)

				// Media Assets
				superAdminGroup.GET("/media", mediaCtrl.AdminGetAll)
				superAdminGroup.POST("/media/upload", mediaCtrl.Upload)
				superAdminGroup.DELETE("/media/:id", mediaCtrl.Delete)

				// Users
				superAdminGroup.GET("/users", userCtrl.GetAll)
				superAdminGroup.POST("/users", userCtrl.Create)
				superAdminGroup.PUT("/users/:id", userCtrl.Update)
				superAdminGroup.DELETE("/users/:id", userCtrl.Delete)

				// Website Settings
				superAdminGroup.GET("/settings", websiteSettingCtrl.GetAll)
				superAdminGroup.PUT("/settings", websiteSettingCtrl.BulkUpsert)
			}
		}
	}

	return r
}
