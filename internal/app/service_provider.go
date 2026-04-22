package app

import (
	"context"
	"github.com/M1steryO/Room-Booking-Service/internal/client"
	"github.com/M1steryO/Room-Booking-Service/internal/client/conference"
	"github.com/M1steryO/Room-Booking-Service/internal/config"
	authrepo "github.com/M1steryO/Room-Booking-Service/internal/repository/auth"
	bookingsrepo "github.com/M1steryO/Room-Booking-Service/internal/repository/bookings"
	roomsrepo "github.com/M1steryO/Room-Booking-Service/internal/repository/rooms"
	schedulesrepo "github.com/M1steryO/Room-Booking-Service/internal/repository/schedules"
	slotsrepo "github.com/M1steryO/Room-Booking-Service/internal/repository/slots"
	authuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/auth"
	bookingsuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/bookings"
	roomsuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/rooms"
	schedulesuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/schedules"
	slotsuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/slots"
	"github.com/M1steryO/Room-Booking-Service/pkg/clock"
	"github.com/M1steryO/Room-Booking-Service/pkg/security"
	"github.com/M1steryO/platform_common/pkg/closer"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type serviceProvider struct {
	httpConfig   config.HTTPConfig
	dbConfig     config.DBConfig
	appConfig    config.AppConfig
	loggerConfig config.LoggerConfig

	pool *pgxpool.Pool

	clock      clock.Clock
	jwtManager *security.JWTManager
	conference client.ConferenceService

	authUC      *authuc.AuthUsecase
	roomsUC     *roomsuc.RoomsUsecase
	schedulesUC *schedulesuc.SchedulesUsecase
	slotsUC     *slotsuc.SlotsUsecase
	bookingsUC  *bookingsuc.BookingsUsecase
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}
func (sp *serviceProvider) LoggerConfig() config.LoggerConfig {
	if sp.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			log.Fatalf("failed to load logger config: %v", err)
		}
		sp.loggerConfig = cfg
	}
	return sp.loggerConfig
}
func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}
		sp.httpConfig = cfg
	}
	return sp.httpConfig
}

func (sp *serviceProvider) DBConfig() config.DBConfig {
	if sp.dbConfig == nil {
		cfg, err := config.NewDBConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		sp.dbConfig = cfg
	}

	return sp.dbConfig
}

func (sp *serviceProvider) AppConfig() config.AppConfig {
	if sp.appConfig == nil {
		cfg, err := config.NewAppConfig()
		if err != nil {
			log.Fatalf("failed to get app config: %v", err)
		}
		sp.appConfig = cfg
	}
	return sp.appConfig
}

func (sp *serviceProvider) Pool(ctx context.Context) *pgxpool.Pool {
	if sp.pool == nil {
		pool, err := pgxpool.Connect(ctx, sp.DBConfig().GetDSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}
		if err = pool.Ping(ctx); err != nil {
			log.Fatalf("failed to ping db: %s", err.Error())
		}
		sp.pool = pool
		closer.Add(func() error {
			pool.Close()
			return nil
		})
	}
	return sp.pool
}

func (sp *serviceProvider) Clock() clock.Clock {
	if sp.clock == nil {
		sp.clock = clock.RealClock{}
	}
	return sp.clock
}

func (sp *serviceProvider) JWTManager() *security.JWTManager {
	if sp.jwtManager == nil {
		sp.jwtManager = security.NewJWTManager(
			sp.AppConfig().JWTSecret(),
			sp.AppConfig().JWTTTL(),
		)
	}
	return sp.jwtManager
}

func (sp *serviceProvider) ConferenceService() client.ConferenceService {
	if sp.conference == nil {
		sp.conference = conference.NewConferenceService(
			sp.AppConfig().ConferenceBaseURL(),
			sp.AppConfig().ConferenceTimeout(),
		)
	}
	return sp.conference
}

func (sp *serviceProvider) AuthUsecase(ctx context.Context) *authuc.AuthUsecase {
	if sp.authUC == nil {
		usersRepo := authrepo.NewUsersRepository(sp.Pool(ctx))
		sp.authUC = authuc.NewAuthUsecase(usersRepo, sp.Clock(), sp.JWTManager())
	}
	return sp.authUC
}

func (sp *serviceProvider) RoomsUsecase(ctx context.Context) *roomsuc.RoomsUsecase {
	if sp.roomsUC == nil {
		roomsRepo := roomsrepo.NewRoomsRepository(sp.Pool(ctx))
		sp.roomsUC = roomsuc.NewRoomsUsecase(roomsRepo)
	}
	return sp.roomsUC
}

func (sp *serviceProvider) SchedulesUsecase(ctx context.Context) *schedulesuc.SchedulesUsecase {
	if sp.schedulesUC == nil {
		roomsRepo := roomsrepo.NewRoomsRepository(sp.Pool(ctx))
		schedulesRepo := schedulesrepo.NewSchedulesRepository(sp.Pool(ctx))
		slotsRepo := slotsrepo.NewSlotsRepository(sp.Pool(ctx))

		sp.schedulesUC = schedulesuc.NewSchedulesUsecase(
			roomsRepo,
			schedulesRepo,
			slotsRepo,
			sp.Clock(),
			sp.AppConfig().SlotHorizonDays(),
		)
	}
	return sp.schedulesUC
}

func (sp *serviceProvider) SlotsUsecase(ctx context.Context) *slotsuc.SlotsUsecase {
	if sp.slotsUC == nil {
		roomsRepo := roomsrepo.NewRoomsRepository(sp.Pool(ctx))
		schedulesRepo := schedulesrepo.NewSchedulesRepository(sp.Pool(ctx))
		slotsRepo := slotsrepo.NewSlotsRepository(sp.Pool(ctx))

		sp.slotsUC = slotsuc.NewSlotsUsecase(
			roomsRepo,
			schedulesRepo,
			slotsRepo,
		)
	}
	return sp.slotsUC
}

func (sp *serviceProvider) BookingsUsecase(ctx context.Context) *bookingsuc.BookingsUsecase {
	if sp.bookingsUC == nil {
		slotsRepo := slotsrepo.NewSlotsRepository(sp.Pool(ctx))
		bookingsRepo := bookingsrepo.NewBookingsRepository(sp.Pool(ctx))

		sp.bookingsUC = bookingsuc.NewBookingsUsecase(
			slotsRepo,
			bookingsRepo,
			sp.ConferenceService(),
			sp.Clock(),
		)
	}
	return sp.bookingsUC
}
