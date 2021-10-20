package sim

import api._
import io.gatling.core.Predef._
import io.gatling.core.structure.ScenarioBuilder
import io.gatling.http.Predef._
import io.gatling.http.protocol.HttpProtocolBuilder

import scala.concurrent.duration._
import scala.util.Random
import scala.math.exp

class BasicSimulation extends Simulation {
  def calcBuzzIndex(twId: Int): Double = 0.01 + 0.99 * exp(twId % 500 / 50.0 - 10.0)

  val httpProtocol: HttpProtocolBuilder = http
    .baseUrl(sys.env("YUMETTER_API") + "/v1") // Here is the root for all relative URLs
    .acceptHeader("application/json;q=0.9,*/*;q=0.8") // Here are the common headers
    .acceptEncodingHeader("gzip, deflate")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  Random.setSeed(System.currentTimeMillis())

  val feeder: Iterator[Map[String, Any]] = Iterator.continually(
    Map(
      "user_name" -> ("U" + System.currentTimeMillis.toString + Thread.currentThread.getId.toString),
      "password" -> System.currentTimeMillis.toString,
    )
  )

  val scn: ScenarioBuilder = scenario("Basic scenario") // A scenario is a chain of requests and pauses
    .feed(feeder)
    .exec(
      Users.post,
      Users.Login.post,
      Users.Me.get,
    )
    .repeat(100)(
      exec(
        Tweets.post,
        repeat(5)(
          Tweets.get.foreach("${tweets}", "tweet") {
            exec{ session =>
              val tweet = session("tweet").as[Map[String, Any]]
              val buzzIndex = calcBuzzIndex(tweet("tw_id").asInstanceOf[Integer])

              session
                .set("tw_id", tweet("tw_id"))
                .set("buzz_index", buzzIndex)
            }
              .doIf(session => Random.nextDouble() < session("buzzIndex").as[Double]) { Tweets.Favourites.put }
              .doIf(session => Random.nextDouble() < session("buzzIndex").as[Double] / 2) { Tweets.postReply }
              .doIf(session => Random.nextDouble() < session("buzzIndex").as[Double] / 1.2) { Tweets.postReply }
          }
        )
      )
    )

  setUp(
    scn.inject(
      incrementUsersPerSec(4)
        .times(50)
        .eachLevelLasting(10 seconds)
        .startingFrom(1)).protocols(httpProtocol),
  )
}
