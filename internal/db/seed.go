package db

import (
	"context"
	"fmt"
	"log"

	"math/rand"

	"github.com/edwrdc/digitally/internal/store"
)

var usernames = []string{
	"james", "emma", "liam", "olivia", "noah", "ava", "william", "sophia", "lucas",
	"isabella", "henry", "mia", "oliver", "charlotte", "alexander", "amelia",
	"benjamin", "harper", "mason", "evelyn", "ethan", "abigail", "daniel",
	"emily", "jacob", "elizabeth", "michael", "sofia", "david", "victoria",
}

var productNames = []string{
	"Professional Video Editing Software", "Digital Art Creation Suite",
	"Music Production DAW", "Cloud Storage Subscription",
	"Website Builder Pro", "Virtual Private Network Service",
	"Password Manager Premium", "Photo Editing Software",
	"Online Course Bundle", "Productivity Suite License",
	"Antivirus Software", "Game Development Engine",
	"3D Modeling Software", "Stock Photo Collection",
	"Audio Book Subscription", "Project Management Tool",
	"Digital Marketing Course", "Programming IDE License",
	"Video Streaming Service", "Digital Sheet Music",
	"Language Learning App", "Meditation App Premium",
	"Graphic Design Templates", "Social Media Management Tool",
	"Digital Art Brushes Pack", "Audio Sample Library",
	"Email Marketing Service", "Digital Magazine Subscription",
	"Online Fitness Program", "Business Analytics Software",
}

var productDescriptions = []string{
	"Unlock the full potential of your video projects with our professional video editing software. Perfect for filmmakers, content creators, and video editors.",
	"Create stunning digital art with our comprehensive suite of tools. Whether you're a hobbyist or a professional, this software has everything you need to bring your ideas to life.",
	"Produce professional-quality music with our Music Production DAW. Record, mix, and master your tracks with ease, and unleash your creativity.",
	"Store and access your files securely with our cloud storage subscription. Keep your important documents, photos, and videos safe and accessible from anywhere.",
	"Build your own website effortlessly with our website builder pro. Choose from a wide range of templates and customize them to fit your brand and vision.",
	"Secure your online activities with our virtual private network service. Protect your data and privacy while browsing the web, and enjoy enhanced online security.",
	"Manage your passwords securely with our password manager premium. Keep your sensitive information safe and easily accessible across all your devices.",
	"Edit and enhance your photos with our photo editing software. Perfect for photographers, graphic designers, and social media managers.",
	"Learn new skills and advance your career with our online course bundle. Access a wide range of courses in various fields, from business to programming.",
	"Boost your productivity with our productivity suite license. Get access to a range of tools and apps that will help you stay organized and efficient.",
	"Protect your devices and data with our antivirus software. Keep your system safe from viruses, malware, and other threats.",
	"Develop your own games with our game development engine. Create your own unique games and unleash your creativity.",
	"Create stunning 3D models with our 3D modeling software. Perfect for architects, designers, and artists.",
	"Access a wide range of high-quality stock photos for your projects. Perfect for designers, photographers, and social media managers.",
	"Subscribe to our audio book subscription and enjoy a wide range of narrated stories and novels. Perfect for readers and listeners.",
	"Manage your projects and tasks with our project management tool. Keep your team organized and on track with our easy-to-use platform.",
	"Learn new languages and improve your skills with our language learning app. Perfect for students, travelers, and language enthusiasts.",
	"Discover a new way to relax and unwind with our meditation app premium. Perfect for anyone looking to improve their mental health and well-being.",
	"Get access to a wide range of graphic design templates to help you create stunning designs for your projects. Perfect for designers and creatives.",
	"Manage your social media accounts with our social media management tool. Perfect for business owners, marketers, and social media managers.",
	"Get access to a wide range of digital art brushes to help you create stunning artwork. Perfect for artists and designers.",
	"Access a wide range of high-quality audio samples to help you create professional-quality music. Perfect for musicians and producers.",
	"Send targeted email campaigns with our email marketing service. Perfect for businesses and marketers looking to grow their audience and increase sales.",
	"Stay up-to-date with the latest news and trends with our digital magazine subscription. Perfect for anyone looking to stay informed and entertained.",
	"Get access to a wide range of online fitness programs to help you stay fit and healthy. Perfect for anyone looking to improve their physical health and well-being.",
	"Analyze your business data with our business analytics software. Perfect for business owners, marketers, and data analysts looking to grow their business and increase sales.",
}

var categories = []string{
	"electronics", "clothing", "home", "garden", "sports", "books", "music",
	"art", "toys", "food", "health", "beauty", "automotive", "travel", "pets",
	"finance", "education", "technology", "gaming", "movies", "music", "sports",
	"art", "toys", "food", "health", "beauty", "automotive", "travel", "pets",
}
var productReviews = []string{
	"This software is amazing! I've been using it for my video editing projects and it's really helped me create professional-quality videos.",
	"I've been using this software for my digital art projects and it's really helped me create stunning artwork. The brushes are amazing!",
	"This software is really easy to use and it's helped me create professional-quality music. I've been using it for my music production projects and it's really helped me create professional-quality music.",
	"Excellent value for money. The features are comprehensive and the interface is intuitive.",
	"Game changer for my business! Customer support is top-notch whenever I need help.",
	"Been using this for 6 months now and it keeps getting better with each update.",
	"The cloud storage is reliable and secure. Never had any issues with file access.",
	"Perfect for beginners but also has advanced features for pros. Highly recommend!",
	"The templates are modern and easy to customize. Saved me so much time.",
	"Great performance and stability. Haven't experienced any crashes or bugs.",
	"The mobile app sync works flawlessly. Love being able to access everything on the go.",
	"Worth every penny! The productivity boost has paid for itself many times over.",
	"Security features are robust. Feel much safer using this for sensitive data.",
	"Learning curve is minimal thanks to the excellent tutorials and documentation.",
	"Regular updates keep adding useful features. Developers clearly listen to feedback.",
	"Integration with other tools is seamless. Works perfectly with my existing workflow.",
	"The community forums are incredibly helpful for tips and troubleshooting.",
	"Fantastic UI/UX design. Everything is exactly where you'd expect it to be.",
	"Processing speed is impressive even with large files.",
	"Cross-platform compatibility works great between my devices.",
	"The analytics features provide valuable insights for my business.",
	"Best in class for its category. Have tried others but keep coming back to this one.",
	"Automation features save me hours every week. Absolutely essential tool.",
}

func Seed(store *store.Storage) error {

	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("error creating user: %v", err)
			return err
		}
	}

	products := generateProducts(200, users)
	for _, product := range products {
		if err := store.Products.Create(ctx, product); err != nil {
			log.Printf("error creating product: %v", err)
			return err
		}
	}

	reviews := generateReviews(200, users, products)
	for _, review := range reviews {
		if err := store.Reviews.Create(ctx, review); err != nil {
			log.Printf("error creating review: %v", err)
			return err
		}
	}

	return nil
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)

	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d@example.com", i),
			Password: "123123",
		}
	}

	return users
}

func generateProducts(n int, users []*store.User) []*store.Product {
	products := make([]*store.Product, n)

	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]

		products[i] = &store.Product{
			UserID:      user.ID,
			Name:        productNames[rand.Intn(len(productNames))],
			Price:       rand.Float64() * 100,
			Description: productDescriptions[rand.Intn(len(productDescriptions))],
			Categories: []string{
				categories[rand.Intn(len(categories))],
				categories[rand.Intn(len(categories))],
			},
		}
	}

	return products
}

func generateReviews(n int, users []*store.User, products []*store.Product) []*store.Review {
	reviews := make([]*store.Review, n)

	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]
		product := products[rand.Intn(len(products))]

		reviews[i] = &store.Review{
			UserID:    user.ID,
			ProductID: product.ID,
			Rating:    rand.Intn(5) + 1,
			Comment:   productReviews[rand.Intn(len(productReviews))],
		}
	}

	log.Println("seeding complete")

	return reviews
}
