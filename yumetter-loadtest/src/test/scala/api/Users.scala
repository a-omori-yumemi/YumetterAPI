package api

import io.gatling.core.Predef._
import io.gatling.core.structure.ChainBuilder
import io.gatling.http.Predef._

object Users {
  val post: ChainBuilder = exec(
    http("creates a new user")
      .post("/users")
      .formParamMap(
        Map(
          "name" -> "${user_name}",
          "password" -> "${password}"
        )
      )
  )

  object Login {
    val post: ChainBuilder = exec(
      http("logs into the service")
        .post("/users/login")
        .formParamMap(
          Map(
            "name" -> "${user_name}",
            "password" -> "${password}"
          )
        )
        .check(status.is(200))
    )
      .exitHereIfFailed
      .exec(
        getCookieValue(
          CookieKey("SESSION")
            .withSecure(true)
            .saveAs("TOO_SECURE_SESSION")
        )
      )
  }

  object Me {
    val get: ChainBuilder = exec(
      http("gets the owner of the session")
        .get("/users/me")
        .check(
          jsonPath("$.usr_id")
            .ofType[Int]
            .saveAs("usr_id")
        )
    )
  }
}
