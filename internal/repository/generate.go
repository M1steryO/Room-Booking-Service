package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/M1steryO/Room-Booking-Service/internal/repository.UserRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/M1steryO/Room-Booking-Service/internal/repository.RoomRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/M1steryO/Room-Booking-Service/internal/repository.ScheduleRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/M1steryO/Room-Booking-Service/internal/repository.SlotRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/M1steryO/Room-Booking-Service/internal/repository.BookingRepository -o ./mocks -s "_minimock.go"
