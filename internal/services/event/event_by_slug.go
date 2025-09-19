package event

import (
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
)

// GetEventBySlug retrieves an event by its slug from the database and writes the response.
//
// Parameters:
//   - wr: An HttpWriter to write the HTTP response
//   - repo: The database repository for event operations
//   - conn: The database connection (not used in this function but included for consistency)
//
// Behavior:
//   - Extracts the 'slug' parameter from the URL
//   - Queries the database for an event matching the slug
//   - If found, responds with the event data and HTTP 200 status
//   - If not found, responds with HTTP 404 status and an error message
//   - If any other error occurs, responds with HTTP 500 status and an error message
func GetEventBySlug(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	slug := utils.GetParam(wr.R, "slug")

	if slug == "" {
		wr.Status(400).Json(
			utils.M{
				"message": "Slug parameter is required",
			},
		)
		return
	}

	event, err := repo.GetEventBySlug(wr.R.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(404).Json(
				utils.M{
					"message": "Event not found",
				},
			)
			return
		}
		wr.Status(500).Json(
			utils.M{
				"message": "Error while fetching event",
				"error":   err.Error(),
			},
		)
		return
	}

	// Fetch associated ticket types
	ticketTypes, err := repo.GetTicketTypesByEventSlug(wr.R.Context(), slug)
	if err != nil {
		wr.Status(500).Json(
			utils.M{
				"message": "Error while fetching ticket types",
				"error":   err.Error(),
			},
		)
		return
	}

	// Fetch associated coupons
	coupons, err := repo.GetCouponsByEventSlug(wr.R.Context(), slug)
	if err != nil {
		wr.Status(500).Json(
			utils.M{
				"message": "Error while fetching coupons",
				"error":   err.Error(),
			},
		)
		return
	}

	// Construct the response with event, ticket types, and coupons
	eventResponse := utils.M{
		"id":                   event.ID,
		"title":                event.Title,
		"description":          event.Description,
		"banner":               event.Banner,
		"icon":                 event.Icon,
		"admin_id":             event.AdminID,
		"start_time":           event.StartTime,
		"end_time":             event.EndTime,
		"location":             event.Location,
		"total_seats":          event.TotalSeats,
		"available_seats":      event.AvailableSeats,
		"slug":                 event.Slug,
		"organizer_name":       event.OrganizerName,
		"organizer_email":      event.OrganizerEmail,
		"organizer_phone":      event.OrganizerPhone,
		"organization":         event.Organization,
		"contact_email":        event.ContactEmail,
		"contact_phone":        event.ContactPhone,
		"refund_policy":        event.RefundPolicy,
		"terms_and_conditions": event.TermsAndConditions,
		"event_type":           event.EventType,
		"category":             event.Category,
		"max_tickets_per_user": event.MaxTicketsPerUser,
		"booking_start_time":   event.BookingStartTime,
		"booking_end_time":     event.BookingEndTime,
		"tags":                 event.Tags,
		"website_url":          event.WebsiteUrl,
		"facebook_url":         event.FacebookUrl,
		"twitter_url":          event.TwitterUrl,
		"instagram_url":        event.InstagramUrl,
		"linkedin_url":         event.LinkedinUrl,
		"venue_name":           event.VenueName,
		"address_line1":        event.AddressLine1,
		"address_line2":        event.AddressLine2,
		"city":                 event.City,
		"state":                event.State,
		"postal_code":          event.PostalCode,
		"country":              event.Country,
		"latitude":             event.Latitude,
		"longitude":            event.Longitude,
		"created_at":           event.CreatedAt,
		"updated_at":           event.UpdatedAt,
		"ticket_types":         ticketTypes,
		"coupons":              coupons,
	}

	wr.Status(200).Json(
		utils.M{
			"message": "Event fetched successfully",
			"event":   eventResponse,
		},
	)
}
