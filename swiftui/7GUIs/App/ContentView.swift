import SwiftUI

struct ContentView: View {
  let views = [
    (
      name: "Counter",
      view: AnyView(Counter())
    ),

    (
      name: "Temperature Converter",
      view: AnyView(TemperatureConverter())
    ),

    //    (name: "Flight Booker", view: Counter()),
    //    (name: "Timer", view: Counter()),
    //    (name: "CRUD", view: Counter()),
    //    (name: "Circle Drawer", view: Counter()),
    //    (name: "Cells", view: Counter()),
  ]

  var body: some View {
    NavigationView {
      List {
        ForEach(views, id: \.name) { (name, view) in
          NavigationLink {
            view
              .navigationBarTitle(name)
              .navigationBarTitleDisplayMode(.inline)
              .padding()
          } label: {
            Text(name)
          }
        }
      }
      .navigationTitle("7 GUIs")

      // this only appears as second columns on iPad and mac.
      views[0].view
    }
  }
}

struct ContentView_Previews: PreviewProvider {
  static var previews: some View {
    Group {
      ContentView()
        .preferredColorScheme(.light)
      ContentView()
        .preferredColorScheme(.dark)
        .previewInterfaceOrientation(.landscapeLeft)
    }
  }
}
