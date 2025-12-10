using System;

namespace Tachograph
{
    /// <summary>
    /// Configures the marshaling process for tachograph files.
    /// </summary>
    public class MarshalOptions
    {
        /// <summary>
        /// Controls whether the marshaler uses RawData fields from parsed messages to reconstruct the file.
        /// 
        /// If true (default), the marshaler will use the RawData fields when available,
        /// applying the "raw data painting" strategy to ensure perfect binary round-tripping
        /// while validating semantic field correctness.
        /// 
        /// If false, the marshaler will always encode from semantic fields, ignoring any RawData fields.
        /// This is useful when semantic fields have been modified and you want to generate new binary data.
        /// </summary>
        public bool UseRawData { get; set; } = true;

        /// <summary>
        /// Serializes a parsed tachograph file into its binary representation.
        /// </summary>
        /// <param name="file">The file to marshal</param>
        /// <returns>Binary data</returns>
        /// <exception cref="ArgumentNullException">Thrown when file is null</exception>
        /// <exception cref="NotSupportedException">Thrown when the file type is not supported</exception>
        public byte[] Marshal(File file)
        {
            if (file == null)
                throw new ArgumentNullException(nameof(file));

            switch (file.Type)
            {
                case FileType.DriverCard:
                    return MarshalDriverCardFile(file.DriverCard);

                case FileType.VehicleUnit:
                    return MarshalVehicleUnitFile(file.VehicleUnit);

                default:
                    throw new NotSupportedException($"Unsupported file type for marshaling: {file.Type}");
            }
        }

        private byte[] MarshalDriverCardFile(DriverCardFile card)
        {
            // TODO: Implement driver card file marshaling
            throw new NotImplementedException("Driver card file marshaling not yet implemented");
        }

        private byte[] MarshalVehicleUnitFile(VehicleUnitFile vu)
        {
            // TODO: Implement VU file marshaling
            throw new NotImplementedException("Vehicle unit file marshaling not yet implemented");
        }
    }
}
