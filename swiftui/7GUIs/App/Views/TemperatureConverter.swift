import SwiftUI

// learn from
// https://github.com/iTwenty/7guis-swiftui/blob/main/Shared/02%20-%20Temperature%20Converter/TempConverterView.swift
private enum ChangeSource {
  case user, indirect
}

struct TemperatureConverter: View {
  @State private var celsius = "0"
  @State private var celsiusSrc = ChangeSource.user
  @State private var fahr = "32"
  @State private var fahrSrc = ChangeSource.user

  var body: some View {
    VStack {
      input("Celsius", $celsius) { str in
        if celsiusSrc == .user, let c = Double(str) {
          fahrSrc = .indirect
          fahr = String(c * 9 / 5 + 32)
        }
        celsiusSrc = .user
      }

      input("Fahrenheit", $fahr) { str in
        if fahrSrc == .user, let f = Double(str) {
          celsiusSrc = .indirect
          celsius = String((f - 32) / 9 * 5)
        }
        fahrSrc = .user
      }

    }
    .padding()
  }

  func input(
    _ label: String,
    _ value: Binding<String>,
    _ changed: @escaping (String) -> Void
  ) -> some View {
    HStack(alignment: .center, spacing: 30) {
      Text(label)
        .frame(maxWidth: .infinity, alignment: .trailing)

      TextField(
        "\(label) value to be converted",
        text: value
      )
      .onChange(of: value.wrappedValue, perform: changed)
      .frame(maxWidth: .infinity)
      .padding()
      .overlay {
        Double(value.wrappedValue) == nil
          ? RoundedRectangle(cornerRadius: 5).stroke(.red, lineWidth: 1)
          : RoundedRectangle(cornerRadius: 5).stroke(.blue, lineWidth: 0)
      }
    }
  }

  func celsiusChange(d: Double) {}
}

struct TemperatureConverter_Previews: PreviewProvider {
  static var previews: some View {
    TemperatureConverter()
  }
}
