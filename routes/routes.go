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
