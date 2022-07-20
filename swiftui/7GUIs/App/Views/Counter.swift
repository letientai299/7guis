import SwiftUI

struct Counter: View {
  @State private var count = 0

  var body: some View {
    VStack {
      Button {
        count += 1
      } label: {
        Text("Count")
      }
      .font(.title)

      Text("\(count)")
        .font(.largeTitle)
    }
  }
}

struct Counter_Previews: PreviewProvider {
  static var previews: some View {
    Counter()
  }
}
