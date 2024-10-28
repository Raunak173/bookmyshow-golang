# Movie Booking Service

This app is a movie booking service where users can browse current movies, select a venue, choose seats, and make reservations. After confirming a booking, users log in or register to proceed. The app manages seat availability, ensuring race conditions are handled during booking, with a 10-minute window to complete the payment before seats are released. Admins can manage venues, movies, showtimes, and bookings. It also has otp verification using twilio.

I had also dockerised it as: https://hub.docker.com/repository/docker/raunak173/bms/general

Also deployed it on AWS EC2 instance at http://ec2-13-232-146-41.ap-south-1.compute.amazonaws.com

The Database used is AWS RDS Postgresql 

Postman Collection: https://www.postman.com/mission-geoscientist-67700153/workspace/raunak-public-apis/collection/28779130-2fc604b8-a41f-4957-b586-ea10091712a7?action=share&creator=28779130

## Customer User Journey

1. **Movie Selection**: 
   - Users can view the movies that are playing for the current week.
   - The user selects a movie.

2. **Venue Selection**: 
   - Based on the movie selected, venues where that movie is playing, along with show timings, are displayed.
   - The user selects a venue and a showtime.

3. **Seat Selection**: 
   - The user is shown the theater layout with available, reserved, and booked seats.
   - The user selects a seat and the number of seats to book.

4. **Booking Confirmation**: 
   - The user clicks on confirm booking.
   - If not logged in, the user is prompted to register or login.
   - User verification is done via OTP on their phone number and email.

5. **Concurrency Management**: 
   - Once confirmed, concurrency conditions are handled to avoid race conditions for seat bookings.
   - The user has 10 minutes to complete the payment; otherwise, the seats are unreserved.

6. **Payment**: 
   - Upon successful payment, the seats are marked as sold.
   - A confirmation email is sent to the user with the booking and ticket details.

---

## Admin User Journey

1. **Admin Login**:
   - The admin logs into the system.

2. **Venue Management**:
   - The admin can add, update, or delete venues.

3. **Movie Management**:
   - The admin can manage which movies are shown at which venues.

4. **Showtime Management**:
   - The admin can add, update, or delete showtimes for different movies at different venues.

5. **Booking Management**:
   - The admin can manage bookings for each venue and showtime.

---

## ER Diagram Relationships

### Users Table

- **Fields**: userId, name, email, password, role (admin or customer)

### Movies Table

- **Fields**: movieId, title, description, releaseDate, duration
- A movie can be associated with multiple venues (Many-to-Many).

### Venues Table

- **Fields**: venueId, name, location
- A venue can host multiple movies (Many-to-Many).

### Movie_Venues Table

- **Fields**: movieId, venueId
- This table establishes the many-to-many relationship between movies and venues.

### ShowTimes Table

- **Fields**: showtimeId, showtime, movieId, venueId
- Each showtime links a specific movie to a specific venue.

### Seats Table

- **Fields**: seatId, seatNumber, isReserved, isBooked, isAvailable, price, showtimeId
- Seats are associated with a specific showtime.

### Orders Table

- **Fields**: orderId, userId, showtimeId, totalPrice
- An order contains booking information for a particular showtime.

### Order_Seats Table

- **Fields**: orderId, seatId
- Links the seats booked in a particular order.

---

## ER Diagram Overview

- **Users** can be either admins or customers.
- **Movies** can be shown at multiple **venues**, and **venues** can host multiple **movies** (Many-to-Many relationship).
- **ShowTimes** link a **movie** to a **venue** at a specific time.
- **Seats** are tied to a specific **showtime**, with details about reservation and booking status.
- **Orders** contain details about the total price and the specific **seats** booked for a **showtime**.

## My Rough Planning Doc

https://docs.google.com/document/d/1V-5DDfVjn_tCE37AQka71sHt51J6D9FnvRBAQxyjd0o/edit?usp=sharing
