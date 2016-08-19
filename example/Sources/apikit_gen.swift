import UIKit
import APIKit

protocol GitHubRequestType: RequestType {
}

extension GitHubRequestType {
    var baseURL: NSURL {
        return NSURL(string: "https://api.github.com")!
    }
}



struct RateLimit {
struct Rate {
let limit: Int
let remaining: Int
let reset: Int
init?(dictionary: [String: AnyObject]) {
    guard let limit = dictionary["limit"] as? Int else {
        return nil
    }
    self.limit = limit
    guard let remaining = dictionary["remaining"] as? Int else {
        return nil
    }
    self.remaining = remaining
    guard let reset = dictionary["reset"] as? Int else {
        return nil
    }
    self.reset = reset

}
}
let rate: Rate?

struct Resources {
struct Core {
let limit: Int
let remaining: Int
let reset: Int
init?(dictionary: [String: AnyObject]) {
    guard let limit = dictionary["limit"] as? Int else {
        return nil
    }
    self.limit = limit
    guard let remaining = dictionary["remaining"] as? Int else {
        return nil
    }
    self.remaining = remaining
    guard let reset = dictionary["reset"] as? Int else {
        return nil
    }
    self.reset = reset

}
}
let core: Core?

struct Search {
let limit: Int
let remaining: Int
let reset: Int
init?(dictionary: [String: AnyObject]) {
    guard let limit = dictionary["limit"] as? Int else {
        return nil
    }
    self.limit = limit
    guard let remaining = dictionary["remaining"] as? Int else {
        return nil
    }
    self.remaining = remaining
    guard let reset = dictionary["reset"] as? Int else {
        return nil
    }
    self.reset = reset

}
}
let search: Search?

init?(dictionary: [String: AnyObject]) {
    guard let core = dictionary["core"] as? [String: AnyObject] else {
        return nil
    }
    self.core = Core(dictionary: core)
    guard let search = dictionary["search"] as? [String: AnyObject] else {
        return nil
    }
    self.search = Search(dictionary: search)

}
}
let resources: Resources?

init?(dictionary: [String: AnyObject]) {
    guard let rate = dictionary["rate"] as? [String: AnyObject] else {
        return nil
    }
    self.rate = Rate(dictionary: rate)
    guard let resources = dictionary["resources"] as? [String: AnyObject] else {
        return nil
    }
    self.resources = Resources(dictionary: resources)

}
}

struct GETRateLimitRequest: GitHubRequestType {
    typealias Response = RateLimit

    var method: HTTPMethod {
        return .GET
    }

    var path: String {
        return "/rate_limit"
    }

    func responseFromObject(object: AnyObject, URLResponse: NSHTTPURLResponse) throws -> Response {
        guard let dictionary = object as? [String: AnyObject], let rateLimit = RateLimit(dictionary: dictionary) else {
            throw ResponseError.UnexpectedObject(object)
        }
        return rateLimit
    }
}

