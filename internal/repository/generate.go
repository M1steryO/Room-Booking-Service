package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/avito-internships/test-backend-1-M1steryO/internal/repository.UserRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/avito-internships/test-backend-1-M1steryO/internal/repository.RoomRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/avito-internships/test-backend-1-M1steryO/internal/repository.ScheduleRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/avito-internships/test-backend-1-M1steryO/internal/repository.SlotRepository -o ./mocks -s "_minimock.go"
//go:generate minimock -i github.com/avito-internships/test-backend-1-M1steryO/internal/repository.BookingRepository -o ./mocks -s "_minimock.go"
