package example

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._
import scala.util.Random
import java.util.Calendar

object WatchTL {
  val watchTL = repeat(Random.nextInt(3), "n")(
    exec(
      http("TL").get("/tweets")
    ).pause(Random.nextInt(3))
  )
}

object RegisterAndLogin {
  val reg = exec{session => 
    val now = Calendar.getInstance()
    val userName = "U" + now.getTimeInMillis().toString() +Thread.currentThread().getId().toString()
    println(userName)
    val password = "pass" + now.getTimeInMillis().toString()
    session.set("user_name", userName).set("password", password)
  }.exec(
    http("Reg").post("/users").formParamMap(Map(
        "name" -> "${user_name}",
        "password" -> "${password}"
    ))
  ).exec(
    http("Login").post("/users/login").formParamMap(Map(
        "name" -> "${user_name}",
        "password" -> "${password}"
    ))
  ).exec(
    getCookieValue(CookieKey("SESSION").withSecure(true).saveAs("TOO_SECURE_SESSION"))
  ).exec(
    addCookie(Cookie("SESSION", "${TOO_SECURE_SESSION}"))
  ).exec(
    http("CheckMe").get("/users/me").check(jsonPath("$..usr_id").ofType[Int].saveAs("usr_id"))
  )
}

object Tweet {
  val tweet = exec(
    http("Tweet").post("/tweets").formParamMap(Map(
      "body" -> ("Hello " + Random.alphanumeric.take(20).mkString)
    ))
  )
}

class BasicItSimulation extends Simulation {

  val httpProtocol = http
    .baseUrl("http://localhost:8000/v1") // Here is the root for all relative URLs
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8") // Here are the common headers
    .acceptEncodingHeader("gzip, deflate")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  val scn = scenario("Basic scenario") // A scenario is a chain of requests and pauses
    .exec(WatchTL.watchTL)
    .exec(RegisterAndLogin.reg)
    .repeat(Random.nextInt(5))(
      exec(Tweet.tweet)
      .exec(WatchTL.watchTL)
    )

  print(2)

  setUp(scn.inject(rampUsers(100).during(30.seconds)).protocols(httpProtocol))
}
