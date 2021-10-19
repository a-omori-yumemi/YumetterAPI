package example

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._
import scala.util.Random
import scala.math.exp
import java.util.Calendar


object WatchTL {
  def calcBuzzIndex(twId: Int) = (0.01 + 0.99 * exp(twId % 500 / 50.0 - 10.0))
  val watchTL = repeat(3, "n")(
    exec(
      http("TL").get("/tweets")
    ).pause(3, 10)
  )

  val randomBin = Iterator.continually( Map( "randomBin" -> Random.nextInt(2)))

  val watchFavReply = repeat(5)(
    exec(
      http("TL").get("/tweets").check(jsonPath("$[*].tweet").ofType[Map[String,Any]].findAll.saveAs("tweets"))
    )
    .foreach("${tweets}", "tweet") {
      exec{session => 
        val tweet = session("tweet").as[Map[String, Any]]
        val buzzIndex = calcBuzzIndex(tweet("tw_id").asInstanceOf[Integer])
        val doFav = Random.nextDouble() < buzzIndex
        val doReply = Random.nextDouble() < buzzIndex/2
        val doLookReply = Random.nextDouble() < buzzIndex/1.2
        session.set("tw_id", tweet("tw_id")).set("doFav", doFav).set("doReply",doReply).set("doLookReply",doLookReply)
      }
      .doIf("${doFav}") {
        exec(
          http("Fav").put("/tweets/${tw_id}/favorites/${usr_id}")
      )}
      .doIf("${doLookReply}") {
        exec(
          http("Show Replies").get("/tweets?replied_to=${tw_id}")
        )
      }
      .doIf("${doReply}") {
        exec(
          http("Reply").post("/tweets").formParamMap{session =>
            Map(
              "body" -> "reply",
              "replied_to" -> session("tw_id").as[Int]
            )
          }
        )
      }
    } 
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
    http("CheckMe").get("/users/me").check(jsonPath("$.usr_id").ofType[Int].saveAs("usr_id"))
  )
}

object Tweet {
  val tweet = exec(
    http("Tweet").post("/tweets").formParamMap(session=> Map(
      "body" -> ("Hello " + Random.alphanumeric.take(20).mkString)
    ))
  )
}

class BasicItSimulation extends Simulation {

  val httpProtocol = http
    .baseUrl(sys.env("YUMETTER_API") + "/v1") // Here is the root for all relative URLs
    .acceptHeader("application/json;q=0.9,*/*;q=0.8") // Here are the common headers
    .acceptEncodingHeader("gzip, deflate")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  val scn = scenario("Basic scenario") // A scenario is a chain of requests and pauses
    .exec(WatchTL.watchTL)
    .exec(RegisterAndLogin.reg)
    .repeat(5)(
      exec(Tweet.tweet)
      .exec(WatchTL.watchFavReply)
    )

  print(2)

  setUp(scn.inject(rampUsers(100).during(30.seconds)).protocols(httpProtocol))
}
